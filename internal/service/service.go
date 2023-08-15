package service

import (
	"context"
	"github.com/stubborn-gaga-0805/prepare2go/pkg"
)

var service *Service

type Service struct {
	ctx context.Context
	pkg *pkg.Pkg

	orderService *OrderService
}

func NewService(ctx context.Context) *Service {
	if service == nil {
		service = &Service{
			ctx: ctx,
			pkg: pkg.NewPkg(),
		}
	}
	return service
}

/*func (s *Service) GetMqProducer() *mq.Producer {
	return s.mq
}*/

// GetService 返回Service实例
func (s *Service) GetService() *Service {
	return s
}

func (s *Service) GetOrderService() *OrderService {
	if s.orderService == nil {
		s.orderService = NewOrderService(s.ctx)
	}
	return s.orderService
}
