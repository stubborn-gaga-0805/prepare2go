package pkg

import (
	"github.com/stubborn-gaga-0805/prepare2go/pkg/jwt"
	"github.com/stubborn-gaga-0805/prepare2go/pkg/logger"
	"github.com/stubborn-gaga-0805/prepare2go/pkg/mysql"
	"github.com/stubborn-gaga-0805/prepare2go/pkg/oss"
	"github.com/stubborn-gaga-0805/prepare2go/pkg/rabbitmq"
	"github.com/stubborn-gaga-0805/prepare2go/pkg/redis"
	"github.com/stubborn-gaga-0805/prepare2go/pkg/wechat"
)

var pkgConfig Config

type Config struct {
	Env     Env           `json:"env" yaml:"env"`
	Logger  logger.Config `json:"logger" yaml:"logger"`
	MQ      MQ            `json:"rabbitmq" yaml:"rabbitmq"`
	Trace   Trace         `json:"trace" yaml:"trace"`
	Oss     oss.Oss       `json:"oss" yaml:"oss"`
	Jwt     jwt.Jwt       `json:"jwt" yaml:"jwt"`
	Wechat  wechat.Config `json:"wechat" yaml:"wechat"`
	TopScrm ScrmTop       `json:"scrmTop" yaml:"scrmTop"`
}

type Env struct {
	AppId      string `json:"appId" yaml:"appId"`
	AppName    string `json:"appName" yaml:"appName"`
	AppVersion string `json:"appVersion" yaml:"appVersion"`
	AppEnv     string `json:"appEnv" yaml:"appEnv"`
}

type Data struct {
	Db     mysql.DB    `json:"db" yaml:"db"`
	Retail mysql.DB    `json:"retail" yaml:"retail"`
	Redis  redis.Redis `json:"redis" yaml:"redis"`
}

type MQ struct {
	RabbitMQ rabbitmq.Config `json:"rabbitMQ"`
	RedisMQ  redis.MQConfig  `json:"redisMQ"`
}

type Trace struct {
	Jaeger Jaeger `json:"jaeger" yaml:"jaeger"`
}

type Jaeger struct {
	EndPoint string `json:"endPoint" yaml:"endPoint"`
}

type ScrmTop struct {
	AppKey          string `json:"appKey" yaml:"appKey"`
	AppSecret       string `json:"appSecret" yaml:"appSecret"`
	ServerUrl       string `json:"serverUrl" yaml:"serverUrl"`
	ConnectTimeount int64  `json:"connectTimeount" yaml:"connectTimeount"`
	ReadTimeout     int64  `json:"readTimeout" yaml:"readTimeout"`
}

func SetConfig(config Config) {
	pkgConfig = config
}

func GetConfig() Config {
	return pkgConfig
}
