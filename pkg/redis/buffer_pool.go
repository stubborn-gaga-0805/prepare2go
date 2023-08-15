package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/fatih/color"
	"github.com/stubborn-gaga-0805/prepare2go/pkg/helpers"
	"github.com/stubborn-gaga-0805/prepare2go/pkg/logger"
	"time"
)

func (sp *senderPool) flushSenderPool(ctx context.Context) {
	if len(sp.pool) == 0 {
		return
	}
	// 量少直接刷
	if len(sp.pool) <= 10 {
		sp.redisClient.LPush(ctx, sp.redisKey, sp.pool...)
		sp.resetSenderPool()
	}
	// 量大分批刷
	var (
		page      = helpers.DivCeil(len(sp.pool), 100)
		blockList = make([][]interface{}, page)
		block     = make([]interface{}, 0, 100)
	)
	for _, msg := range sp.pool {
		if len(block) < cap(block) {
			block = append(block, msg)
		} else {
			blockList = append(blockList, block)
			block = make([]interface{}, 0, 100)
		}
	}
	// 管道发送
	pip := sp.redisClient.Pipeline()
	defer pip.Close()

	for _, messages := range blockList {
		pip.LPush(ctx, sp.redisKey, messages...)
	}
	if _, err := pip.Exec(ctx); err != nil {
		logger.Helper().Errorf("exec redis pip error... %v", err)
		return
	}
	sp.resetSenderPool()
	return
}

func (sp *senderPool) resetSenderPool() {
	if len(sp.pool) > 0 {
		sp.pool = make([]interface{}, 0, sp.poolCapacity)
	}
	return
}

func (sp *senderPool) poolCleaner(ctx context.Context) {
	ticker := time.NewTicker(sp.flushTicker)
	defer ticker.Stop()

	fmt.Printf("%s (status: %d/%d) [id: %s]--->[redisKey: %s] \n", color.BlueString("%s", "RedisMQ Producer Buffer Pool Running..."), sp.poolSize, sp.poolCapacity, sp.id, sp.redisKey)
	for {
		select {
		case <-ticker.C:
			sp.punchDuration()
			sp.flushSenderPool(ctx)
		}
	}
}

func (sp *senderPool) punchDuration() {
	sp.running = time.Since(sp.createdAt)
}

func (sp *senderPool) terminated(ctx context.Context) {
	sp.running = time.Since(sp.createdAt)
	sp.stopAt = time.Now()
	sp.flushSenderPool(ctx)
	fmt.Printf("%s (running: %s) [id: %s]--->[redisKey: %s] \n", color.YellowString("%s", "RedisMQ Producer Buffer Pool Stopped..."), sp.running.String(), sp.id, sp.redisKey)
}

func (rp *receivePool) messageKeeper(ctx context.Context) {
	result, err := rp.redisClient.BRPop(ctx, 0, rp.redisKey).Result()
	if err != nil {
		logger.Helper().Errorf("rp.redisClient.BRPop err: %v, redisKey: %s", err, rp.redisKey)
		return
	}
	for _, msg := range result {
		var msgStruct messagePackage
		if err = json.Unmarshal([]byte(msg), &msgStruct); err != nil {
			logger.Helper().Errorf("json.Unmarshal Error. err: %v", err)
			continue
		}
		rp.pool <- &msgStruct
	}
	return
}

func (rp *receivePool) punchDuration() {
	rp.running = time.Since(rp.createdAt)
}

func (rp *receivePool) terminated(ctx context.Context) {
	rp.running = time.Since(rp.createdAt)
	rp.stopAt = time.Now()
	rp.flushReceivePool(ctx)
	fmt.Printf("%s (running: %s) [id: %s]--->[redisKey: %s] \n", color.YellowString("%s", "RedisMQ Producer Buffer Pool Stopped..."), rp.running.String(), rp.id, rp.redisKey)
}

func (rp *receivePool) flushReceivePool(ctx context.Context) {
	var (
		residueNum = len(rp.pool)
		msgList    = make([]interface{}, 0)
	)
	if residueNum > 0 {
		for msg := range rp.pool {
			msgList = append(msgList, msg)
		}
	}
	rp.redisClient.RPush(ctx, rp.redisKey, msgList...)
	return
}

func (rp *receivePool) consumerHandler() {
	select {
	case <-rp.stopChan:
		logger.Helper().Infof("[%s] consumer goroutine stopped...[duration: %s]\n", time.Now().Format(time.RFC3339), time.Since(rp.createdAt))
		return
	default:
		for msg := range rp.pool {
			rp.handleFunc(helpers.GetContextWithRequestId(), msg)
		}
	}
	return
}
