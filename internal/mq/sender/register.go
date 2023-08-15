package sender

import (
	"github.com/stubborn-gaga-0805/prepare2go/internal/consts"
	"github.com/stubborn-gaga-0805/prepare2go/internal/mq/contruct"
	"github.com/stubborn-gaga-0805/prepare2go/pkg/rabbitmq"
	"github.com/stubborn-gaga-0805/prepare2go/pkg/redis"
	"time"
)

// 定义所有生产者
var (
	OrderProducer  *rabbitmq.Producer // 订单
	MemberProducer *redis.MqProducer
)

type Producer struct {
	AMQP    *rabbitmq.AMQP
	RedisMQ *redis.MessageQueue
}

// register 注册生产者
func (p *Producer) register() *Producer {
	// RedisMQ Producer
	OrderProducer = rabbitmq.NewProducer(p.AMQP, rabbitmq.ProducerConfig{
		Type:       rabbitmq.Direct,
		Exchange:   contruct.ExchangeGoFrameOrder,
		RoutingKey: contruct.RoutingGoFrameOrder,
	}).Start()

	// RedisMQ Producer
	MemberProducer = redis.NewProducer(p.RedisMQ, redis.ProducerConfig{
		RedisKey:              consts.MemberInfoMQKey,
		BufferPoolFlushTicker: time.Second * 2,
		BufferCapacity:        1e3,
		WarningSize:           1e4,
	}).Start()

	return p
}

func (p *Producer) Run() *Producer {
	reg := p.register()
	reg.AMQP.Running()
	reg.RedisMQ.Running()
	return p
}
