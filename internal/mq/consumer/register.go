package consumer

import (
	"github.com/stubborn-gaga-0805/prepare2go/internal/mq/consumer/handle"
	"github.com/stubborn-gaga-0805/prepare2go/internal/mq/contruct"
	"github.com/stubborn-gaga-0805/prepare2go/pkg/logger"
	"github.com/stubborn-gaga-0805/prepare2go/pkg/rabbitmq"
)

type OrderMsg struct {
	Id      int64  `json:"id"`
	Tid     int64  `json:"tid"`
	Cid     int64  `json:"cid"`
	Sid     int64  `json:"sid"`
	OrderSn string `json:"order_sn"`
}

type Consumer struct {
	AMQP   *rabbitmq.AMQP
	Logger *logger.Logger
	Handle *handle.MqHandle
	// 声明消费者
	OrderConsumer *rabbitmq.Consumer // 订单
}

// register 注册消费者
func (c *Consumer) register() *Consumer {
	c.OrderConsumer = rabbitmq.NewConsumer(c.AMQP, rabbitmq.ConsumerConfig{
		Exchange:   contruct.ExchangeGoFrameOrder,
		RoutingKey: contruct.RoutingGoFrameOrder,
		Queue:      contruct.QueueGoFrameOrder,
		AutoAck:    true,
		ResStruct:  OrderMsg{},
		Handel:     c.Handle.TestOrderSyncConsumer,
	}).SetName("order-consumer").Start()

	return c
}

func (c *Consumer) Run() *Consumer {
	c.register().AMQP.StartConsumer()
	return c
}
