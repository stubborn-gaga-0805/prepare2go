package contruct

import "github.com/stubborn-gaga-0805/prepare2go/pkg/rabbitmq"

// 定义所有的 Exchanges
const (
	ExchangeGoFrameOrder  rabbitmq.Exchange = "ExchangeGoFrameOrder"
	ExchangeGoFrameMember rabbitmq.Exchange = "ExchangeGoFrameMember"
)

// 定义所有的 RoutingKey
const (
	RoutingGoFrameOrder  rabbitmq.RoutingKey = "RoutingGoFrameOrder"
	RoutingGoFrameMember rabbitmq.RoutingKey = "RoutingGoFrameMember"
)

// 定义所有的 Queue
const (
	QueueGoFrameOrder  rabbitmq.Queue = "QueueGoFrameOrder"
	QueueGoFrameMember rabbitmq.Queue = "QueueGoFrameMember"
)
