package rocket_mq

// 消息类型的常量定义
const (
	MsgTypeNormal      MessageType = 8 << iota // 普通消息
	MsgTypeFIFO                                // 顺序消息
	MsgTypeDelay                               // 延迟消息
	MsgTypeTransaction                         // 事务消息
)

const (
	MsgTypeNormalName      = "Normal"      // 普通消息
	MsgTypeFIFOName        = "FIFO"        // 顺序消息
	MsgTypeDelayName       = "Delay"       // 延迟消息
	MsgTypeTransactionName = "Transaction" // 事务消息
)

const (
	producerIdKey    = "ProducerId"
	requestIdKey     = "RequestId"
	userMessageIdKey = "messageId"
	messageTypeKey   = "messageType"
)
