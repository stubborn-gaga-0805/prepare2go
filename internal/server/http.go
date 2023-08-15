package server

import (
	"context"
	"fmt"
	"github.com/google/wire"
	"github.com/kataras/iris/v12"
	recover2 "github.com/kataras/iris/v12/middleware/recover"
	"github.com/kataras/iris/v12/middleware/requestid"
	"github.com/stubborn-gaga-0805/prepare2go/api/router"
	"github.com/stubborn-gaga-0805/prepare2go/configs/conf"
	"github.com/stubborn-gaga-0805/prepare2go/internal/job"
	"github.com/stubborn-gaga-0805/prepare2go/pkg/helpers"
	"time"
)

// ProviderSet is server providers.
var ProviderSet = wire.NewSet(NewHttpServer)

type Server struct {
	*iris.Application
	closers []func()
	conf    conf.Http

	job    *job.Job
	router *router.Handler
}

// NewHttpServer 示例化Http服务
func NewHttpServer() *Server {
	var cfg = conf.GetConfig()

	app := iris.New().SetName(cfg.Env.AppName)
	app.Use(recover2.New())
	app.Use(requestid.New())

	hs := &Server{
		Application: app,
		router:      router.NewRouter(helpers.GetContextWithRequestId(), app),
		conf:        cfg.Server.Http,
	}
	if cfg.WithCronJob {
		hs.job = job.NewJob(context.Background()).RegisterCustomJobs()
	}

	return hs
}

func (hs *Server) Start() error {
	hs.ConfigureHost(func(su *iris.Supervisor) {
		// Set timeouts. More than enough, normally we use 20-30 seconds.
		su.Server.ReadTimeout = 5 * time.Minute
		su.Server.WriteTimeout = 5 * time.Minute
		su.Server.IdleTimeout = 10 * time.Minute
		su.Server.ReadHeaderTimeout = 2 * time.Minute
	})

	addr := fmt.Sprintf("%s:%d", hs.conf.Addr, hs.conf.Port)
	fmt.Println(fmt.Sprintf("HttpServer started successfully, listening address: %s\n\n", addr))

	return hs.Listen(addr)
}

func (hs *Server) CrontabSwitchOn() bool {
	return hs.job != nil
}

func (hs *Server) StartCrontab() {
	hs.job.Start()
}

func (hs *Server) StopCrontab() {
	hs.job.Stop()
}
