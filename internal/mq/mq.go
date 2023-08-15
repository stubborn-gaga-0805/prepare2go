package mq

import (
	"context"
	"github.com/stubborn-gaga-0805/prepare2go/configs/conf"
	"github.com/stubborn-gaga-0805/prepare2go/internal/mq/consumer"
	"github.com/stubborn-gaga-0805/prepare2go/internal/mq/consumer/handle"
	"github.com/stubborn-gaga-0805/prepare2go/internal/mq/sender"
	"github.com/stubborn-gaga-0805/prepare2go/pkg/rabbitmq"
	"github.com/stubborn-gaga-0805/prepare2go/pkg/redis"
)

type MQ struct {
	p *sender.Producer
	c *consumer.Consumer
}

func NewMQConsumer(ctx context.Context) {
	var cfg = conf.GetConfig()
	if !cfg.WithOutMQ {
		new(MQ).newConsumer(ctx).Run()
	}
}

func NewMQProducer(ctx context.Context) *sender.Producer {
	var cfg = conf.GetConfig()
	if cfg.WithOutMQ {
		return nil
	}
	return new(MQ).newProducer(ctx).Run()
}

func (m *MQ) newProducer(ctx context.Context) *sender.Producer {
	var cfg = conf.GetConfig()
	return &sender.Producer{
		AMQP:    rabbitmq.NewAMQP(ctx, cfg.MQ.RabbitMQ),
		RedisMQ: redis.NewMessageQueue(ctx, cfg.MQ.RedisMQ),
	}
}

func (m *MQ) newConsumer(ctx context.Context) *consumer.Consumer {
	var cfg = conf.GetConfig()
	return &consumer.Consumer{
		AMQP:   rabbitmq.NewAMQP(ctx, cfg.MQ.RabbitMQ),
		Handle: handle.NewMqHandle(ctx),
	}
}
