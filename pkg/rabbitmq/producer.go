package rabbitmq

import (
	"context"
	"encoding/json"
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/samber/lo"
	"github.com/stubborn-gaga-0805/prepare2go/pkg/consts"
	"github.com/stubborn-gaga-0805/prepare2go/pkg/helpers"
	"github.com/stubborn-gaga-0805/prepare2go/pkg/logger"
	"reflect"
	"time"
)

const (
	Direct ConnType = 16 << iota
	FanOut
	Topic
	Headers
)

const (
	requestIdKey = "requestId"
)

type ConnType uint8

type Producer struct {
	producerKey  string
	exchange     string
	routingKey   string
	producerType string

	ap       *AMQP
	conn     *amqp.Connection
	messages chan []byte
}

type ProducerConfig struct {
	Exchange   Exchange
	RoutingKey RoutingKey
	Type       ConnType
}

var (
	legalConnType   = []ConnType{Direct, FanOut, Topic, Headers}
	connTypeMapping = map[ConnType]string{
		Direct:  "direct",
		FanOut:  "fanout",
		Topic:   "topic",
		Headers: "headers",
	}
)

func NewProducer(ap *AMQP, c ProducerConfig) *Producer {
	checkProducerConfig(c)
	return &Producer{
		ap:           ap,
		producerType: connTypeMapping[c.Type],
		exchange:     string(c.Exchange),
		routingKey:   string(c.RoutingKey),
		messages:     make(chan []byte, 3),
	}
}

func (p *Producer) Start() *Producer {
	p.producerKey = p.key()
	p.ap.addProducer(p)
	return p
}

func checkProducerConfig(cfg ProducerConfig) {
	if len(cfg.Exchange) == 0 {
		panic("NewProducer Err: Error Exchange Name")
	}
	if len(cfg.RoutingKey) == 0 {
		panic("NewProducer Err: Error RoutingKey")
	}
	if !lo.Contains(legalConnType, cfg.Type) {
		panic("NewProducer Err: illegal Type")
	}
}

func (p *Producer) running(ctx context.Context) (err error) {
	p.conn = p.ap.conn()
	channel, err := p.conn.Channel()
	if err != nil {
		return err
	}
	if err = channel.ExchangeDeclarePassive(p.exchange, p.producerType, true, true, false, false, nil); err != nil {
		channel, err := p.conn.Channel()
		if err != nil {
			return err
		}
		if err = channel.ExchangeDeclare(
			p.exchange,
			p.producerType,
			true,
			true,
			false,
			false,
			nil,
		); err != nil {
			return err
		}
	}

	defer p.conn.Close()
	defer channel.Close()

	return nil
}

func (p *Producer) Push(ctx context.Context, message any) {
	var (
		requestId = ctx.Value(consts.ContextRequestIdKey)
		messageId = helpers.GenUUID()
	)
	producer, ok := p.ap.producerList[p.producerKey]
	if !ok {
		panic("Producer Not Fount!")
	}
	producer.conn = producer.ap.conn()
	channel, err := producer.conn.Channel()
	if err != nil {
		panic(err)
	}
	err = channel.PublishWithContext(
		helpers.GetContextWithRequestId(),
		p.exchange,
		p.routingKey,
		false,
		false,
		amqp.Publishing{
			ContentType:     "application/json",
			ContentEncoding: "UTF-8",
			Timestamp:       time.Now(),
			MessageId:       messageId,
			Body:            producer.parsePushMessage(message),
			Headers: map[string]interface{}{
				requestIdKey: requestId,
			},
		},
	)
	if err != nil {
		logger.HelperWithContext(ctx).Errorf("PublishWithContext Err: %v", err)
		return
	}
	logger.HelperWithContext(ctx).Infof("Producer发送消息[%s]--->[exchange: %s] to [routingKey: %s]", messageId, p.exchange, p.routingKey)
	return
}

func (p *Producer) parsePushMessage(message interface{}) []byte {
	var (
		t   = reflect.TypeOf(message)
		msg []byte
		err error
	)
	switch t.Kind() {
	case reflect.String:
		msg = []byte(fmt.Sprintf("%s", message))
	case reflect.Map:
		fallthrough
	case reflect.Slice:
		fallthrough
	case reflect.Array:
		fallthrough
	case reflect.Struct:
		if msg, err = json.Marshal(message); err != nil {
			panic(err)
		}
	default:
		panic("不支持的消息类型！")
	}
	return msg
}

func (p *Producer) key() string {
	return helpers.Md5Encrypt(fmt.Sprintf("%s+%s+%s", p.producerType, p.exchange, p.routingKey))
}
