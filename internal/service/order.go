package service

import (
	"context"
	"github.com/stubborn-gaga-0805/prepare2go/internal/mq/sender"
	"github.com/stubborn-gaga-0805/prepare2go/pkg/helpers"
	"github.com/stubborn-gaga-0805/prepare2go/pkg/logger"
)

var orderService *OrderService

type OrderMsg struct {
	Id      int64  `json:"id"`
	Tid     int64  `json:"tid"`
	Cid     int64  `json:"cid"`
	Sid     int64  `json:"sid"`
	OrderSn string `json:"order_sn"`
}

type OrderService struct {
	*Service
}

func NewOrderService(ctx context.Context) *OrderService {
	if orderService == nil {
		orderService = &OrderService{
			NewService(ctx),
		}
	}
	return orderService
}

func (s *OrderService) TestProducerMq(ctx context.Context, opt ...string) error {
	logger.Helper().Infof("Handle(func TestProducerMq())消息内容[%+v]", opt[0])
	sender.OrderProducer.Push(ctx, OrderMsg{1, 2, 3, 4, helpers.GenSerialNo(20)})
	return nil
}
