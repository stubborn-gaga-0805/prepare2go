package handle

import (
	"context"
	"github.com/stubborn-gaga-0805/prepare2go/internal/job"
	"github.com/stubborn-gaga-0805/prepare2go/internal/service"
)

type MqHandle struct {
	*service.Service
	*job.Job
}

func NewMqHandle(ctx context.Context) *MqHandle {
	return &MqHandle{
		service.NewService(ctx),
		job.NewJob(ctx),
	}
}
