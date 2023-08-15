package controller

import (
	"context"
)

type Controller struct {
	OrderController *OrderController
}

func NewControllerHandler(ctx context.Context) *Controller {
	return &Controller{
		NewOrderController(ctx),
	}
}
