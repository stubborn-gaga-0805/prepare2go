package rocket_mq

import (
	"context"
	"fmt"
	"github.com/apache/rocketmq-clients/golang/v5/credentials"
	"github.com/fatih/color"
	"github.com/stubborn-gaga-0805/prepare2go/pkg/helpers"
	"github.com/stubborn-gaga-0805/prepare2go/pkg/logger"
	"time"
)

type ConsumerConfig struct {
	Name              string
	ConsumerGroup     string
	Topics            []string
	DialDuration      time.Duration
	InvisibleDuration time.Duration
	MaxMessageNum     int
	Handler           ConsumerHandler
}

type Consumer struct {
	Id                string
	Name              string
	Handler           ConsumerHandler
	startAt           time.Time
	duration          time.Duration
	invisibleDuration time.Duration
	maxMessageNum     int32
	instance          golang.SimpleConsumer
}

type ConsumerHandler func(ctx context.Context, message Message) error

func (mq *RocketMq) NewConsumer(cfg ConsumerConfig) *Consumer {
	defer func() {
		if err := recover(); err != nil {
			logger.Helper().Errorf("New RocketMQ Consumer[%s] Failed... err: %v", cfg.Name, err)
		}
	}()

	instance, err := golang.NewSimpleConsumer(&golang.Config{
		Endpoint:      mq.cfg.EndPoint,
		ConsumerGroup: cfg.ConsumerGroup,
		Credentials: &credentials.SessionCredentials{
			AccessKey:    mq.cfg.AccessKey,
			AccessSecret: mq.cfg.AccessSecret,
		},
	}, golang.WithAwaitDuration(cfg.DialDuration),
	)
	if err != nil {
		panic(nil)
	}
	consumer := &Consumer{
		Id:                helpers.GenUUID(),
		Name:              cfg.Name,
		instance:          instance,
		invisibleDuration: cfg.InvisibleDuration,
		maxMessageNum:     int32(cfg.MaxMessageNum),
	}
	mq.addConsumer(consumer)

	return consumer
}

func (consumer *Consumer) Start() error {
	defer func() {
		if err := recover(); err != nil {
			logger.Helper().Errorf("Start RocketMQ Consumer[%s] Failed... err: %v", consumer.Name, err)
		}
	}()

	if err := consumer.instance.Start(); err != nil {
		panic(err)
	}
	go consumer.receiver()

	return nil
}

func (consumer *Consumer) Shutdown() error {
	consumer.duration = time.Since(consumer.startAt)
	if err := consumer.instance.GracefulStop(); err != nil {
		fmt.Printf("%s (running: %s) [id: %s]--->[name: %s] \n", color.RedString("%s", "Stop RocketMQ Consumer Failed..."), consumer.duration.String(), consumer.Id, consumer.Name)
	}
	fmt.Printf("%s (running: %s) [id: %s]--->[name: %s] \n", color.YellowString("%s", "RocketMQ Consumer Stopped..."), consumer.duration.String(), consumer.Id, consumer.Name)
	return nil
}

func (consumer *Consumer) receiver() {
	consumer.startAt = time.Now()
	fmt.Printf("RocketMQ Consumer[%s][%s] Running...\n", consumer.Name, consumer.startAt.Format(time.DateTime))
	for {
		receive, err := consumer.instance.Receive(context.TODO(), consumer.maxMessageNum, consumer.invisibleDuration)
		if err != nil {
			return
		}
		for _, msg := range receive {
			var (
				message    = msg
				ctx        = helpers.GetContextWithRequestId()
				properties = msg.GetProperties()
			)
			if _, ok := properties[requestIdKey]; !ok {
				properties[requestIdKey] = helpers.GenUUID()
			}
			if _, ok := properties[userMessageIdKey]; !ok {
				properties[userMessageIdKey] = ""
			}
			if _, ok := properties[messageTypeKey]; !ok {
				properties[messageTypeKey] = ""
			}
			messagePackage := Message{
				Body:          msg.GetBody(),
				MsgId:         msg.GetMessageId(),
				RequestId:     properties[requestIdKey],
				UserMessageId: properties[userMessageIdKey],
				MessageGroup:  fmt.Sprintf("%v", msg.GetMessageGroup()),
			}
			messagePackage.MessageType = messagePackage.GetMessageType(properties[messageTypeKey])
			go func() {
				if err := consumer.Handler(ctx, messagePackage); err != nil {
					logger.Helper().Errorf("consumer message failed: %v. messageId: %s", err, message.GetMessageId())
					return
				}
				// 消费成功Ack
				if err := consumer.instance.Ack(ctx, message); err != nil {
					logger.Helper().Errorf("message ack failed: %v. messageId: %s", err, message.GetMessageId())
				}
				return
			}()
		}
	}
}
