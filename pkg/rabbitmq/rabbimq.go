package rabbitmq

import (
	"context"
	"fmt"
	"github.com/fatih/color"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/stubborn-gaga-0805/prepare2go/pkg/helpers"
	"github.com/stubborn-gaga-0805/prepare2go/pkg/logger"
	"github.com/stubborn-gaga-0805/prepare2go/pkg/wework"
	"net"
	"reflect"
	"runtime/debug"
	"sync"
	"time"
)

type (
	Exchange   string
	RoutingKey string
	Queue      string
)

type AMQP struct {
	ctx      context.Context
	cfg      Config
	endpoint string

	producerList map[string]*Producer
	consumerList map[string]*Consumer
}

type msgPackage struct {
	msgType    string
	msgContent []byte
}

type ReceiveMsgPackage struct {
	RequestId  string
	MsgId      string
	SendTime   time.Time
	msgContent interface{}
	originMsg  []byte

	resType reflect.Type
}

func NewAMQP(ctx context.Context, cfg Config) *AMQP {
	return &AMQP{
		ctx:          ctx,
		cfg:          cfg,
		endpoint:     fmt.Sprintf("amqp://%s:%s@%s:%s/%s", cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.Vhost),
		producerList: make(map[string]*Producer, 0),
		consumerList: make(map[string]*Consumer, 0),
	}
}

func (a *AMQP) Running() {
	for _, producer := range a.producerList {
		go func() {
			err := producer.running(a.ctx)
			if err != nil {
				logger.HelperWithContext(a.ctx).Errorf("StartProducer Error. Err: %v, sender: %+v", err, producer)
				panic("StartProducer Error")
			}
		}()
		fmt.Printf("%s [Exchange: %s]--->[RoutingKey: %s]=>[Type: %s] \n", color.BlueString("%s", "RabbitMQ Producer Running..."), producer.exchange, producer.routingKey, producer.producerType)
	}
	return
}

func (a *AMQP) StartConsumer() {
	var wg = new(sync.WaitGroup)
	for _, consumer := range a.consumerList {
		wg.Add(1)
		go consumer.running(wg)
		fmt.Printf("Consumer[name: %s] Running... [Exchange: %s]--->[RoutingKey: %s, queue: %s] \n", consumer.name, consumer.exchange, consumer.routingKey, consumer.queue)
	}
	wg.Wait()
	return
}

func (a *AMQP) addProducer(pd *Producer) {
	a.producerList[pd.producerKey] = pd
	return
}

func (a *AMQP) addConsumer(cs *Consumer) {
	a.consumerList[cs.consumerKey] = cs
	return
}

func (a *AMQP) conn() *amqp.Connection {
	defer func() {
		if err := recover(); err != nil {
			stack := debug.Stack()
			logger.Helper().Errorf("conn rabbitMQ failed: %v", err)
			wework.PanicBroadcast(err, helpers.GetRequestIdFromContext(a.ctx), string(stack))
		}
	}()

	var (
		err    error
		conn   *amqp.Connection
		config = amqp.Config{
			Dial: func(network, addr string) (net.Conn, error) {
				return net.DialTimeout(network, addr, a.cfg.MaxDialTimeout)
			},
		}
	)
	if conn, err = amqp.DialConfig(a.endpoint, config); err != nil {
		panic(err)
	}
	return conn
}
