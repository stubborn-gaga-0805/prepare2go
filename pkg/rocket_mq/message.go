package rocket_mq

import (
	"errors"
	"fmt"
	"github.com/apache/rocketmq-clients/golang/v5"
	"github.com/samber/lo"
	"time"
)

var (
	legalMessageTypes  = []MessageType{MsgTypeNormal, MsgTypeFIFO, MsgTypeDelay, MsgTypeTransaction}
	messageTypeMapping = map[MessageType]string{
		MsgTypeNormal:      MsgTypeNormalName,
		MsgTypeFIFO:        MsgTypeFIFOName,
		MsgTypeDelay:       MsgTypeDelayName,
		MsgTypeTransaction: MsgTypeTransactionName,
	}
	messageNameTypeMapping = map[string]MessageType{
		MsgTypeNormalName:      MsgTypeNormal,
		MsgTypeFIFOName:        MsgTypeFIFO,
		MsgTypeDelayName:       MsgTypeDelay,
		MsgTypeTransactionName: MsgTypeTransaction,
	}
)

type TransactionCheckFunc func(msg *golang.MessageView) golang.TransactionResolution

type MessageType int

type Message struct {
	RequestId     string
	MsgId         string
	UserMessageId string
	MessageType   MessageType
	Body          any
	Topic         string
	MessageGroup  string
	DelayTime     time.Duration
}

func (m *Message) IsLegalType() error {
	if !(lo.Contains(legalMessageTypes, m.MessageType)) {
		return errors.New(fmt.Sprintf("Illeage message type: %v", m.MessageType))
	}
	return nil
}

func (m *Message) IsLegalTopic() error {
	if len(m.Topic) == 0 {
		return errors.New(fmt.Sprintf("Illeage topic: %s", m.Topic))
	}
	return nil
}

func (m *Message) GetMessageType(messageTypeName string) MessageType {
	if _, ok := messageNameTypeMapping[messageTypeName]; !ok {
		return MsgTypeNormal
	}
	return messageNameTypeMapping[messageTypeName]
}

func (messageType MessageType) ToString() string {
	if _, ok := messageTypeMapping[messageType]; !ok {
		return ""
	}
	return messageTypeMapping[messageType]
}
