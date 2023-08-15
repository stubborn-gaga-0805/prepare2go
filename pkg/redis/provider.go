package redis

import (
	"context"
	"github.com/go-redis/redis/v8"
	"github.com/stubborn-gaga-0805/prepare2go/pkg/helpers"
	"github.com/stubborn-gaga-0805/prepare2go/pkg/logger"
	"github.com/stubborn-gaga-0805/prepare2go/pkg/wework"
	"runtime/debug"
)

type Instance struct {
	client *redis.Client

	messageQueue *MessageQueue
}

func NewRedis(ctx context.Context, rc Redis) *Instance {
	defer func() {
		if err := recover(); err != nil {
			stack := debug.Stack()
			logger.Helper().Errorf("init redis failed: %v", err)
			wework.PanicBroadcast(err, helpers.GetRequestIdFromContext(ctx), string(stack))
		}
	}()
	client, err := New(ctx, rc)
	if err != nil {
		panic(err)
	}

	return &Instance{
		client: client,
	}
}

func (instance *Instance) Client() *redis.Client {
	return instance.client
}
