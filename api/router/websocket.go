package router

import (
	"fmt"
	socketio "github.com/googollee/go-socket.io"
	"github.com/stubborn-gaga-0805/prepare2go/internal/service"
	"github.com/stubborn-gaga-0805/prepare2go/pkg"
	"github.com/stubborn-gaga-0805/prepare2go/pkg/helpers"
	"github.com/stubborn-gaga-0805/prepare2go/pkg/logger"
)

type WsRouter struct {
	*service.Service
	*pkg.Pkg
}

func NewWsRouter() *WsRouter {
	return &WsRouter{
		service.NewService(helpers.GetContextWithRequestId()),
		pkg.NewPkg(),
	}
}

func (r *WsRouter) RegisterHandle(server *socketio.Server) *socketio.Server {
	defer func() {
		if err := recover(); err != nil {
			logger.Helper().Errorf("panic: %v", err)
		}
	}()
	server.OnConnect("/ws", func(s socketio.Conn) error {
		s.SetContext(helpers.GenUUID())
		fmt.Printf("[ws] Connected: %s, requestID: %s\n", s.ID(), fmt.Sprintf("%s", s.Context()))
		return nil
	})
	server.OnEvent("/ws", "message", func(s socketio.Conn, jsonMsg interface{}) {
		ctx := helpers.GetContextWithRequestId()
		s.SetContext(helpers.GetRequestIdFromContext(ctx))

		_, err := helpers.WsJsonMsg2Map(jsonMsg)
		if err != nil {
			s.Emit("error", err.Error())
		}
		s.Emit("reply", "hello")
		return
	})
	server.OnEvent("/ws", "shutdown", func(s socketio.Conn) string {
		last := s.Context().(string)
		s.Emit("reply", last)
		s.Close()
		return last
	})
	server.OnError("/ws", func(s socketio.Conn, e error) {
		fmt.Printf("[ws] Error: %v\n", e)
	})
	server.OnDisconnect("/", func(s socketio.Conn, reason string) {
		fmt.Printf("[ws] Closed. %s\n", reason)
	})
	return server
}
