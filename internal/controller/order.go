package controller

import (
	"context"
	"errors"
	"github.com/kataras/iris/v12"
	"github.com/stubborn-gaga-0805/prepare2go/api/response"
	"github.com/stubborn-gaga-0805/prepare2go/internal/service"
)

var orderController *OrderController

type OrderController struct {
	orderService *service.OrderService
}

func NewOrderController(ctx context.Context) *OrderController {
	if orderController == nil {
		orderController = &OrderController{
			service.NewOrderService(ctx),
		}
	}
	return orderController
}

func (c *OrderController) SendMsg(ctx iris.Context) {
	err := c.orderService.TestProducerMq(ctx, "Hello")
	if err != nil {
		response.ErrorJson(ctx, errors.New("mq error"))
		return
	}
	response.SuccessJson(ctx)
}
