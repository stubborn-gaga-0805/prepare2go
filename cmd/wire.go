//go:build wireinject
// +build wireinject

package cmd

import (
	"context"
	"github.com/google/wire"
	"github.com/stubborn-gaga-0805/prepare2go/internal/server"
)

func initServer(context.Context) (*server.Server, func(), error) {
	panic(wire.Build(server.ProviderSet))
}

func initWebSocket(context.Context) (*server.Websocket, func(), error) {
	panic(wire.Build(server.WsProviderSet))
}
