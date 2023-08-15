package mysql

import (
	"context"
	"github.com/stubborn-gaga-0805/prepare2go/pkg/helpers"
	"github.com/stubborn-gaga-0805/prepare2go/pkg/logger"
	"github.com/stubborn-gaga-0805/prepare2go/pkg/wework"
	"gorm.io/gorm"
	"runtime/debug"
)

func NewMysql(ctx context.Context, c DB) *gorm.DB {
	defer func() {
		if err := recover(); err != nil {
			stack := debug.Stack()
			logger.Helper().Errorf("Init Mysql Error...%v", err)
			wework.PanicBroadcast(err, helpers.GetRequestIdFromContext(ctx), string(stack))
		}
	}()

	db, err := New(ctx, c)
	if err != nil {
		panic(err)
	}

	return db
}
