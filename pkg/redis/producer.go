package redis

import (
	"context"
	"fmt"
	"github.com/fatih/color"
	"github.com/stubborn-gaga-0805/prepare2go/pkg/helpers"
	"github.com/stubborn-gaga-0805/prepare2go/pkg/logger"
	"time"
)

const (
	DefaultProducerWarningNum int = 1e4
)

type MqProducer struct {
	producerId      string
	RedisKey        string
	warningSize     int
	startAt         time.Time
	runningDuration time.Duration
	mq              *MessageQueue
	withSenderPool  bool
}

type ProducerConfig struct {
	RedisKey              string
	BufferPoolFlushTicker time.Duration
	BufferCapacity        int
	WarningSize           int
}

type ProducerStatus struct {
	Id             string
	RedisKey       string
	BufferCapacity int
	BufferSize     int
}

func NewProducer(mq *MessageQueue, cfg ProducerConfig) *MqProducer {
	if len(cfg.RedisKey) == 0 || mq == nil {
		panic("init redis mq producer failed...!")
	}
	p := &MqProducer{
		producerId:     helpers.GenUUID(),
		mq:             mq,
		RedisKey:       cfg.RedisKey,
		warningSize:    DefaultProducerWarningNum,
		withSenderPool: false,
	}
	if cfg.BufferCapacity > 0 {
		mq.createSenderPool(p.Id(), cfg)
		p.withSenderPool = true
	}
	if cfg.BufferPoolFlushTicker.Seconds() == 0 {
		cfg.BufferPoolFlushTicker = time.Second * 3
	}
	if cfg.WarningSize > 0 {
		p.warningSize = cfg.WarningSize
	}
	return p
}

func (producer *MqProducer) Start() *MqProducer {
	producer.addProducer()
	return producer
}

func (producer *MqProducer) running() {
	producer.startAt = time.Now()
	producer.runningDuration = time.Since(producer.startAt)
	fmt.Printf("%s [id: %s]--->[redisKey: %s] \n", color.BlueString("%s", "RedisMQ Producer Running..."), producer.Id(), producer.RedisKey)
	pool, ok := producer.mq.senderPools[producer.Id()]
	if ok {
		// SenderBuffer 监听
		go pool.poolCleaner(producer.mq.ctx)
	}
}

func (producer *MqProducer) Id() string {
	return producer.producerId
}

func (producer *MqProducer) Push(ctx context.Context, message any) {
	var (
		requestId = helpers.GetRequestIdFromContext(ctx)
		messageId = helpers.GenUUID()
		msg       = messagePackage{
			requestId: requestId,
			messageId: messageId,
			Body:      message,
		}
	)
	if producer.withSenderPool {
		if err := producer.mq.dropToSenderPool(ctx, producer.producerId, msg); err != nil {
			logger.Helper().Errorf("producer.mq.dropToSenderPool err: %v", err)
		}
	} else {
		if err := producer.dropToRedis(ctx, msg); err != nil {
			logger.Helper().Errorf("producer.mq.client.LPush err: %v", err)
		}
	}
	return
}

func (producer *MqProducer) dropToRedis(ctx context.Context, msg messagePackage) error {
	return producer.mq.client.LPush(ctx, producer.RedisKey, msg).Err()
}

func (producer *MqProducer) addProducer() {
	_, ok := producer.mq.producerList[producer.RedisKey]
	if ok {
		fmt.Printf("list: %+v\n", producer.mq.producerList)
		panic(fmt.Sprintf("duplicated redis producer key[%s]...", producer.RedisKey))

	}
	producer.mq.producerList[producer.RedisKey] = producer
	return
}
