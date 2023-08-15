package redis

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/stubborn-gaga-0805/prepare2go/pkg/helpers"
	"github.com/stubborn-gaga-0805/prepare2go/pkg/logger"
	"time"
)

type messagePackage struct {
	requestId  string
	messageId  string
	producerId string
	Body       interface{}
}

type MessageQueue struct {
	ctx    context.Context
	client *redis.Client

	producerList map[string]*MqProducer
	consumerList map[string]*MqConsumer

	senderPools  map[string]*senderPool
	receivePools map[string]*receivePool
}

type senderPool struct {
	id           string
	redisKey     string
	poolCapacity int
	poolSize     int
	pool         []interface{}
	createdAt    time.Time
	stopAt       time.Time
	running      time.Duration
	flushTicker  time.Duration
	redisClient  *redis.Client
}

type receivePool struct {
	id           string
	redisKey     string
	poolCapacity int
	poolSize     int
	pool         chan *messagePackage
	createdAt    time.Time
	stopAt       time.Time
	running      time.Duration
	redisClient  *redis.Client
	handleFunc   ConsumerHandel
	stopChan     chan struct{}
}

func NewMessageQueue(ctx context.Context, cfg MQConfig) *MessageQueue {
	redisCfg := Redis{
		Addr:         cfg.Addr,
		Username:     cfg.Username,
		Password:     cfg.Password,
		Db:           cfg.Db,
		MaxRetries:   cfg.MaxRetries,
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
		DialTimeout:  cfg.DialTimeout,
		PoolSize:     cfg.PoolSize,
		PoolTimeout:  cfg.PoolTimeout,
	}
	logger.WithRequestId(helpers.GetRequestIdFromContext(ctx))
	return &MessageQueue{
		ctx:          ctx,
		client:       NewRedis(ctx, redisCfg).Client(),
		producerList: make(map[string]*MqProducer, 0),
		consumerList: make(map[string]*MqConsumer, 0),
		senderPools:  make(map[string]*senderPool, 0),
		receivePools: make(map[string]*receivePool, 0),
	}
}

func (mq *MessageQueue) Running() {
	for _, producer := range mq.producerList {
		producer.running()
	}
	return
}

func (mq *MessageQueue) Stop() {
	for _, sp := range mq.senderPools {
		sp.terminated(mq.ctx)
	}
	for _, rp := range mq.receivePools {
		rp.terminated(mq.ctx)
	}
	return
}

func (mq *MessageQueue) createSenderPool(id string, cfg ProducerConfig) {
	mq.senderPools[id] = &senderPool{
		id:           id,
		redisKey:     cfg.RedisKey,
		poolCapacity: cfg.BufferCapacity,
		pool:         make([]interface{}, 0, cfg.BufferCapacity),
		createdAt:    time.Now(),
		flushTicker:  cfg.BufferPoolFlushTicker,
		redisClient:  mq.client,
	}
	return
}

func (mq *MessageQueue) dropToSenderPool(ctx context.Context, id string, msg messagePackage) error {
	pool, ok := mq.senderPools[id]
	if !ok {
		return errors.New(fmt.Sprintf("Sender Buffer Pool (%s) Not Found!", id))
	}
	// 缓冲池中的数据刷到redis队列中
	if len(pool.pool) >= pool.poolCapacity {
		pool.flushSenderPool(ctx)
	}
	msgStr, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	pool.pool = append(pool.pool, string(msgStr))

	return nil
}

func (mq *MessageQueue) createReceivePool(id string, cfg ConsumerConfig) {
	mq.receivePools[id] = &receivePool{
		id:           id,
		redisKey:     cfg.RedisKey,
		poolCapacity: cfg.BufferCapacity,
		poolSize:     0,
		pool:         make(chan *messagePackage, cfg.BufferCapacity),
		createdAt:    time.Now(),
		redisClient:  mq.client,
	}
	return
}
