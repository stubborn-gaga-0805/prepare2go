package repo

import (
	"context"
	"github.com/go-redis/redis/v8"
	"github.com/stubborn-gaga-0805/prepare2go/configs/conf"
	pkgmysql "github.com/stubborn-gaga-0805/prepare2go/pkg/mysql"
	pkgredis "github.com/stubborn-gaga-0805/prepare2go/pkg/redis"
	"gorm.io/gorm"
)

var repo *Repo

type Repo struct {
	db       *gorm.DB
	redisCli redis.Cmdable
}

func NewRepo(ctx context.Context) *Repo {
	var cfg = conf.GetConfig()
	if repo == nil {
		repo = &Repo{
			db:       pkgmysql.NewMysql(ctx, cfg.Data.Db),
			redisCli: pkgredis.NewRedis(ctx, cfg.Data.Redis).Client(),
		}
	}

	return repo
}

func ReloadDB(ctx context.Context) {
	var cfg = conf.GetConfig()
	if repo == nil {
		repo = NewRepo(ctx)
	}
	repo.db = pkgmysql.NewMysql(ctx, cfg.Data.Db)

	return
}

func ReloadRedis(ctx context.Context) {
	var cfg = conf.GetConfig()
	if repo == nil {
		repo = NewRepo(ctx)
	}
	repo.redisCli = pkgredis.NewRedis(ctx, cfg.Data.Redis).Client()

	return
}
