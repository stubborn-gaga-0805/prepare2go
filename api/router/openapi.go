package router

import (
	"fmt"
	"github.com/kataras/iris/v12"
	"time"
)

func (h *Handler) registerOpenApiRoutes(app *iris.Application) {
	// Ping
	app.Get("/ping", func(ctx iris.Context) {
		ctx.WriteString(fmt.Sprintf("[%s] %s", time.Now().Format("2006-01-02 15:04:05"), "...pong!"))
		ctx.Done()
	})

	app.Get("/mq/send", h.Controller.OrderController.SendMsg)
}
