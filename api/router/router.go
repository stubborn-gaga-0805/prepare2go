package router

import (
	"context"
	"github.com/kataras/iris/v12"
	"github.com/stubborn-gaga-0805/prepare2go/api/middleware"
	"github.com/stubborn-gaga-0805/prepare2go/internal/controller"
)

type Handler struct {
	*controller.Controller
	m *middleware.Handler
}

// NewRouter 初始化路由配置
func NewRouter(ctx context.Context, app *iris.Application) *Handler {
	r := &Handler{
		controller.NewControllerHandler(ctx),
		middleware.NewMiddleware(ctx),
	}
	// 公用中间件
	app.Use(r.m.Wrapping)
	// 注册路由
	r.registerOpenApiRoutes(app)

	return r
}
