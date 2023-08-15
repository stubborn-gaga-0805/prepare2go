package rocket_mq

import (
	"fmt"
	"github.com/apache/rocketmq-clients/golang/v5/credentials"
	"github.com/stubborn-gaga-0805/prepare2go/pkg/logger"
)

type RocketMq struct {
	cfg    Config
	client golang.Client

	ProducerList map[string]*Producer
	ConsumerList map[string]*Consumer
}

func NewClient(cfg Config) *RocketMq {
	defer func() {
		if err := recover(); err != nil {
			logger.Helper().Errorf("New RocketMQ Client Failed... err: %v", err)
		}
	}()

	client, err := golang.NewClient(&golang.Config{
		Endpoint:      cfg.EndPoint,
		NameSpace:     cfg.Region,
		ConsumerGroup: cfg.Group,
		Credentials: &credentials.SessionCredentials{
			AccessKey:    cfg.AccessKey,
			AccessSecret: cfg.AccessSecret,
		}})
	if err != nil {
		panic(err)
	}
	return &RocketMq{
		cfg:          cfg,
		client:       client,
		ProducerList: make(map[string]*Producer, 0),
	}
}

func (mq *RocketMq) Id() string {
	return mq.client.GetClientID()
}

func (mq *RocketMq) Stop() error {
	// shutdown all producers graceful
	for _, producer := range mq.ProducerList {
		if err := producer.Shutdown(); err != nil {
			logger.Helper().Errorf("Stop Producer[%s][%s] Failed... Err: %v", producer.Name, producer.Id, err)
			continue
		}
		fmt.Printf("Producer[%s][%s] Stopped...\n", producer.Name, producer.Id)
	}
	// shutdown all consumers graceful
	for _, consumer := range mq.ConsumerList {
		if err := consumer.Shutdown(); err != nil {
			logger.Helper().Errorf("Stop Consumer[%s][%s] Failed... Err: %v", consumer.Name, consumer.Id, err)
			continue
		}
		fmt.Printf("Consumer[%s][%s] Stopped...\n", consumer.Name, consumer.Id)
	}
	fmt.Println("RocketMQ Client Stopped...")

	return mq.client.GracefulStop()
}

func (mq *RocketMq) Running() {
	for _, producer := range mq.ProducerList {
		if err := producer.Start(); err != nil {
			logger.Helper().Errorf("Start Producer[%s][%s] Failed... Err: %v", producer.Name, producer.Id, err)
			continue
		}
		fmt.Printf("Producer[%s][%s] Running...\n", producer.Name, producer.Id)
	}
}

func (mq *RocketMq) addProducer(producer *Producer) {
	mq.ProducerList[producer.Id] = producer
}

func (mq *RocketMq) addConsumer(consumer *Consumer) {
	mq.ConsumerList[consumer.Id] = consumer
}
