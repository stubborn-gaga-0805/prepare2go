package rabbitmq

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/stubborn-gaga-0805/prepare2go/pkg/consts"
	"github.com/stubborn-gaga-0805/prepare2go/pkg/helpers"
	"github.com/stubborn-gaga-0805/prepare2go/pkg/logger"
	"go.uber.org/zap"
	"reflect"
	"sync"
)

type ConsumerHandel func(ctx context.Context, message *ReceiveMsgPackage)

type Consumer struct {
	consumerKey string

	ap   *AMQP
	conn *amqp.Connection

	name       string
	exchange   string
	routingKey string
	queue      string
	autoACK    bool
	resType    reflect.Type
	handel     ConsumerHandel
}

type ConsumerConfig struct {
	Exchange   Exchange
	RoutingKey RoutingKey
	Queue      Queue
	ResStruct  any
	AutoAck    bool
	Handel     ConsumerHandel
}

func NewConsumer(ap *AMQP, cfg ConsumerConfig) *Consumer {
	checkConsumerConfig(cfg)
	c := &Consumer{
		ap:         ap,
		exchange:   string(cfg.Exchange),
		routingKey: string(cfg.RoutingKey),
		queue:      string(cfg.Queue),
		autoACK:    cfg.AutoAck,
		handel:     cfg.Handel,
	}
	c.resType = c.reflectReceiveStruct(cfg.ResStruct)
	return c
}

func (c *Consumer) SetName(name string) *Consumer {
	c.name = name
	return c
}

func (c *Consumer) Start() *Consumer {
	c.consumerKey = c.key()
	c.ap.addConsumer(c)
	return c
}

func (c *Consumer) running(wg *sync.WaitGroup) {
	defer wg.Done()

	c.conn = c.ap.conn()
	channel, err := c.conn.Channel()
	if err != nil {
		panic(err)
	}
	defer c.conn.Close()

	c.queueExisted(channel)
	channel, err = c.conn.Channel()
	if err != nil {
		panic(err)
	}
	defer channel.Close()

	if err = channel.QueueBind(
		c.queue,
		c.routingKey,
		c.exchange,
		false,
		nil,
	); err != nil {
		panic(err)
	}
	// 消费消息
	messages, err := channel.Consume(
		c.queue,
		c.name,
		c.autoACK,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		panic(err)
	}
	for message := range messages {
		var (
			receive = new(ReceiveMsgPackage)
			msg     = reflect.New(c.resType).Interface()
		)
		if err = json.Unmarshal(message.Body, msg); err != nil {
			logger.Helper().Errorf("json.Unmarshal Error. %v", err)
			continue
		}
		receive.RequestId = helpers.IntToString(message.Headers[requestIdKey])
		receive.originMsg = message.Body
		receive.MsgId = message.MessageId
		receive.msgContent = msg
		receive.SendTime = message.Timestamp
		receive.resType = c.resType

		if len(receive.RequestId) == 0 {
			receive.RequestId = helpers.GenUUID()
		}
		logger.Helper().Infof("[exchange: %s]-[routingKey: %s]-[queue: %s]--->Consumer收到消息[%s]", c.exchange, c.routingKey, c.queue, receive.MsgId)
		c.handel(context.WithValue(c.ap.ctx, consts.ContextRequestIdKey, receive.RequestId), receive)
	}
	return
}

func (c *Consumer) queueExisted(channel *amqp.Channel) {
	if _, err := channel.QueueDeclarePassive(c.queue, true, false, false, false, nil); err != nil {
		channel, err := c.conn.Channel()
		if err != nil {
			panic(err)
		}
		if _, err := channel.QueueDeclare(c.queue, true, false, false, false, nil); err != nil {
			panic(err)
		}
	}
	return
}

func (c *Consumer) key() string {
	return helpers.Md5Encrypt(fmt.Sprintf("%s+%s+%s+%s", c.name, c.exchange, c.routingKey, c.queue))
}

func checkConsumerConfig(cfg ConsumerConfig) {
	if len(cfg.Exchange) == 0 {
		panic("NewConsumer Err: Error Exchange Name")
	}
	if len(cfg.RoutingKey) == 0 {
		panic("NewConsumer Err: Error RoutingKey")
	}
	if len(cfg.Queue) == 0 {
		panic("NewConsumer Err: Error Queue")
	}
	if cfg.Handel == nil {
		panic("NewConsumer Err: Undefined Handel")
	}
}

func (c *Consumer) reflectReceiveStruct(res interface{}) reflect.Type {
	return reflect.TypeOf(res)
}

func (res *ReceiveMsgPackage) Content(ptr interface{}) error {
	if err := json.Unmarshal(res.originMsg, ptr); err != nil {
		zap.S().Errorf("json.Unmarshal err: %v", err)
		return errors.New(fmt.Sprintf("解析消息结构体错误, msgId: %s, requestId: %s", res.MsgId, res.RequestId))
	}
	res.msgContent = ptr
	return nil
}
