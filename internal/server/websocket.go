package server

import (
	"fmt"
	"github.com/google/wire"
	socketio "github.com/googollee/go-socket.io"
	"github.com/gorilla/websocket"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/context"
	recover2 "github.com/kataras/iris/v12/middleware/recover"
	"github.com/kataras/iris/v12/middleware/requestid"
	"github.com/stubborn-gaga-0805/prepare2go/api/router"
	"github.com/stubborn-gaga-0805/prepare2go/configs/conf"
	"net/http"
)

var WsProviderSet = wire.NewSet(NewWebsocketServer)

type Websocket struct {
	*iris.Application

	conf   conf.WS
	router *router.WsRouter
}

func NewWebsocketServer() *Websocket {
	var cfg = conf.GetConfig()

	app := iris.New().SetName(cfg.Env.AppName)
	app.Use(recover2.New())
	app.Use(requestid.New())

	ws := &Websocket{
		Application: app,
		conf:        cfg.Server.WS,
		router:      router.NewWsRouter(),
	}

	return ws
}

func (ws *Websocket) Start() (err error) {
	var upGrader = websocket.Upgrader{
		ReadBufferSize:  1 << 10,
		WriteBufferSize: 1 << 10,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	ws.Get("/ws", func(ctx *context.Context) {
		conn, err := upGrader.Upgrade(ctx.ResponseWriter(), ctx.Request(), nil)
		if err != nil {
			fmt.Println(fmt.Sprintf("WebSocket server failed to start %v\n", err))
			return
		}
		defer conn.Close()

		for {
			// 读取客户端发送的消息
			messageType, p, err := conn.ReadMessage()
			if err != nil {
				fmt.Println(fmt.Sprintf("WebSocket: Failed to read message %v\n", err))
				break
			}

			fmt.Println(fmt.Sprintf("WebSocket Received the message: %s type: %d\n", string(p), messageType))

			// 发送消息给客户端
			err = conn.WriteMessage(messageType, []byte("ok"))
			if err != nil {
				fmt.Println(fmt.Sprintf("WebSocket: Failed to send message %v\n", err))
				break
			}
		}
	})

	addr := fmt.Sprintf("%s:%d", ws.conf.Addr, ws.conf.Port)
	fmt.Println(fmt.Sprintf("WebSocket started successfully, listening address: %s\n", addr))
	if err = ws.Listen(addr, iris.WithoutPathCorrection); err != nil {
		fmt.Println(fmt.Sprintf("websocket err: %s\n\n", err))
		return err
	}

	return nil
}

func (ws *Websocket) StartWithSocketIO() (err error) {
	wsServer := ws.router.RegisterHandle(socketio.NewServer(nil))
	go func() {
		err := wsServer.Serve()
		if err != nil {
			panic(err)
		}
	}()
	defer wsServer.Close()
	ws.HandleMany("GET POST", "/socket.io/{any:path}", iris.FromStd(wsServer))

	addr := fmt.Sprintf("%s:%d", ws.conf.Addr, ws.conf.Port)
	fmt.Println(fmt.Sprintf("WebSocket started successfully, listening address: %s\n", addr))
	if err = ws.Listen(addr, iris.WithoutPathCorrection); err != nil {
		fmt.Println(fmt.Sprintf("websocket err: %s\n\n", err))
		return err
	}
	return nil
}
