package secretariat

import (
	"context"
	"fmt"
	"github.com/stubborn-gaga-0805/prepare2go/pkg/helpers"
	"os/signal"
	"syscall"
	"time"
)

type Secretary struct {
	Name     string
	workflow *Workflow

	id              string
	ctx             context.Context
	startAt         time.Time
	endAt           time.Time
	runningDuration time.Duration
}

type SecretaryConfig struct {
	Name        string
	BufferSize  int
	ScannerRate time.Duration
	WithCompass bool
}

func NewSecretary(ctx context.Context, cfg SecretaryConfig) *Secretary {
	return &Secretary{
		ctx:      ctx,
		id:       helpers.GenNanoId(10),
		Name:     cfg.Name,
		workflow: NewWorkflow(ctx, cfg.BufferSize, cfg.ScannerRate, cfg.WithCompass),
	}
}

func (secretary *Secretary) Put(msg interface{}) {
	secretary.workflow.Put(msg)
}

func (secretary *Secretary) Start() {
	signalCtx, signalStop := signal.NotifyContext(secretary.ctx, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)
	defer signalStop()

	secretary.workflow.running()
	secretary.startAt = time.Now()
	select {
	case <-signalCtx.Done():
		secretary.endAt = time.Now()
		secretary.workflow.stop()
		fmt.Printf("[%s] Secretary[%s] Done...[duration: %s]\n", time.Now().Format("2006-01-02T15:04:05.999999"), secretary.Name, time.Since(secretary.startAt).String())
	}
}
