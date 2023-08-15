package redis

import (
	"context"
	"github.com/stubborn-gaga-0805/prepare2go/pkg/helpers"
	"time"
)

const (
	DefaultConsumerBufferCapacity int = 1e3
)

type ConsumerHandel func(ctx context.Context, message *messagePackage)

type MqConsumer struct {
	ConsumerId      string
	RedisKey        string
	startAt         time.Time
	runningDuration time.Duration
	goroutineNum    int
	autoAck         bool
	mq              *MessageQueue
}

type ConsumerConfig struct {
	RedisKey       string
	BufferCapacity int
	GoroutineNum   int
	AutoAck        bool
}

type ConsumerStatus struct {
	Id             string
	RedisKey       string
	BufferCapacity int
	BufferSize     int
}

func NewConsumer(mq *MessageQueue, cfg ConsumerConfig) *MqConsumer {
	if len(cfg.RedisKey) == 0 || mq == nil {
		panic("init redis mq consumer failed...!")
	}
	p := &MqConsumer{
		ConsumerId:   helpers.GenUUID(),
		mq:           mq,
		RedisKey:     cfg.RedisKey,
		goroutineNum: cfg.GoroutineNum,
		autoAck:      cfg.AutoAck,
	}
	if cfg.BufferCapacity == 0 {
		cfg.BufferCapacity = DefaultConsumerBufferCapacity
	}
	mq.createReceivePool(p.ConsumerId, cfg)

	return p
}
