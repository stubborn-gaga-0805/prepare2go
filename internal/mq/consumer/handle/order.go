package handle

import (
	"context"
	"fmt"
	"github.com/stubborn-gaga-0805/prepare2go/pkg/consts"
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

func (h *MqHandle) TestOrderSyncConsumer(ctx context.Context, msg *rabbitmq.ReceiveMsgPackage) {
	logger.WithRequestId(fmt.Sprintf("%s", ctx.Value(consts.ContextRequestIdKey)))
	var OrderMsg = &OrderMsg{}
	if err := msg.Content(OrderMsg); err != nil {

		return
	}
	fmt.Printf("OrderMsg: %#v, sn: %s", OrderMsg, OrderMsg.OrderSn)
	return
}
