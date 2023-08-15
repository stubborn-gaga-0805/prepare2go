package rocket_mq

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/apache/rocketmq-clients/golang/v5/credentials"
	"github.com/fatih/color"
	"github.com/stubborn-gaga-0805/prepare2go/pkg/helpers"
	"github.com/stubborn-gaga-0805/prepare2go/pkg/logger"
	"reflect"
	"time"
)

type ProducerConfig struct {
	Name               string
	Topics             []string
	TransactionChecker TransactionCheckFunc
}

type Producer struct {
	Id                 string
	Name               string
	topics             []string
	startAt            time.Time
	duration           time.Duration
	instance           golang.Producer
	transactionChecker *golang.TransactionChecker
}

func (mq *RocketMq) NewProducer(cfg *ProducerConfig) *Producer {
	defer func() {
		if err := recover(); err != nil {
			logger.Helper().Errorf("New RocketMQ Producer[%s] Failed... err: %v", cfg.Name, err)
		}
	}()
	transactionChecker := &golang.TransactionChecker{
		Check: cfg.TransactionChecker,
	}
	instance, err := golang.NewProducer(&golang.Config{
		Endpoint: mq.cfg.EndPoint,
		Credentials: &credentials.SessionCredentials{
			AccessKey:    mq.cfg.AccessKey,
			AccessSecret: mq.cfg.AccessSecret,
		},
	},
		golang.WithTopics(cfg.Topics...),
		golang.WithTransactionChecker(transactionChecker),
	)
	if err != nil {
		panic(nil)
	}

	producer := &Producer{
		Id:                 helpers.GenUUID(),
		Name:               cfg.Name,
		instance:           instance,
		topics:             cfg.Topics,
		transactionChecker: transactionChecker,
	}
	mq.addProducer(producer)

	return producer
}

func (producer *Producer) Start() error {
	defer func() {
		if err := recover(); err != nil {
			logger.Helper().Errorf("Start RocketMQ Producer[%s] Failed... err: %v", producer.Name, err)
		}
	}()

	if err := producer.instance.Start(); err != nil {
		panic(err)
	}
	producer.startAt = time.Now()
	fmt.Printf("RocketMQ Producer[%s][%s] Running...\n", producer.Name, producer.startAt.Format(time.DateTime))

	return nil
}

func (producer *Producer) Shutdown() error {
	producer.duration = time.Since(producer.startAt)
	if err := producer.instance.GracefulStop(); err != nil {
		fmt.Printf("%s (running: %s) [id: %s]--->[name: %s] \n", color.RedString("%s", "Stop RocketMQ Producer Failed..."), producer.duration.String(), producer.Id, producer.Name)
	}
	fmt.Printf("%s (running: %s) [id: %s]--->[name: %s] \n", color.YellowString("%s", "RocketMQ Producer Stopped..."), producer.duration.String(), producer.Id, producer.Name)
	return nil
}

func (producer *Producer) Push(ctx context.Context, msgPackage Message) error {
	// 校验消息参数
	if err := msgPackage.IsLegalType(); err != nil {
		return err
	}
	if err := msgPackage.IsLegalTopic(); err != nil {
		return err
	}
	// 构建rocketMQ消息结构体
	msg := &golang.Message{
		Topic: msgPackage.Topic,
		Body:  producer.parseMessage(msgPackage.Body),
	}
	msg.SetTag(producer.Name)
	msg.AddProperty(producerIdKey, producer.Id)
	msg.AddProperty(requestIdKey, helpers.GetRequestIdFromContext(ctx))
	msg.AddProperty(userMessageIdKey, helpers.GenUUID())
	msg.AddProperty(messageTypeKey, msgPackage.MessageType.ToString())
	// 根据消息类型
	switch msgPackage.MessageType {
	case MsgTypeNormal:
		break
	case MsgTypeFIFO:
		if len(msgPackage.MessageGroup) == 0 {
			return errors.New("FIFO Msg `MessageGroup` required")
		}
		msg.SetMessageGroup(msgPackage.MessageGroup)
	case MsgTypeDelay:
		if msgPackage.DelayTime.Seconds() == 0 {
			return errors.New("delay Msg `DelayTime` required")
		}
		delayTime := time.Now().Add(msgPackage.DelayTime)
		if delayTime.IsZero() {
			return errors.New("illegal `DelayTime`")
		}
		msg.SetDelayTimestamp(delayTime)
	case MsgTypeTransaction:
		if producer.transactionChecker.Check == nil {
			return errors.New("transaction Msg `CheckFunc` required")
		}
	}
	// 事务消息需要特殊处理
	if msgPackage.MessageType == MsgTypeTransaction {
		transaction := producer.instance.BeginTransaction()
		resp, err := producer.instance.SendWithTransaction(ctx, msg, transaction)
		if err != nil {
			return err
		}
		for _, r := range resp {
			logger.Helper().Infof("Producer[%s], Sender Resp: [%#v], OriginMsg: [%#v]", producer.Name, r, msgPackage)
		}
		return transaction.Commit()
	} else {
		resp, err := producer.instance.Send(ctx, msg)
		if err != nil {
			return err
		}
		for _, r := range resp {
			logger.Helper().Infof("Producer[%s], Sender Resp: [%#v], OriginMsg: [%#v]", producer.Name, r, msgPackage)
		}
	}
	return nil
}

func (producer *Producer) parseMessage(message interface{}) []byte {
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
