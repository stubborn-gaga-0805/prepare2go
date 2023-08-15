package conf

import (
	"github.com/stubborn-gaga-0805/prepare2go/pkg/mysql"
	"github.com/stubborn-gaga-0805/prepare2go/pkg/redis"
)

type Data struct {
	Db    mysql.DB    `json:"db" yaml:"db"`
	Redis redis.Redis `json:"redis" yaml:"redis"`
}
