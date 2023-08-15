package job

import (
	"context"
	"github.com/stubborn-gaga-0805/prepare2go/internal/mq/sender"
	"github.com/stubborn-gaga-0805/prepare2go/pkg/helpers"
)

func RedisMqProduce(ctx context.Context, opt ...string) error {
	total := helpers.StringToInt(opt[0])
	for i := 1; i <= total; i++ {
		sender.MemberProducer.Push(ctx, i+1)
	}
	return nil
}
