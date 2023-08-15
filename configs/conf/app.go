package conf

import (
	"fmt"
	"github.com/spf13/viper"
	"github.com/stubborn-gaga-0805/prepare2go/pkg"
	"github.com/stubborn-gaga-0805/prepare2go/pkg/jwt"
	"github.com/stubborn-gaga-0805/prepare2go/pkg/logger"
	"github.com/stubborn-gaga-0805/prepare2go/pkg/oss"
	"github.com/stubborn-gaga-0805/prepare2go/pkg/wechat"
)

type App struct {
	Server Server        `json:"server" yaml:"server"`
	Env    pkg.Env       `json:"env" yaml:"env"`
	Data   Data          `json:"data" yaml:"data"`
	Logger logger.Config `json:"logger" yaml:"logger"`
	MQ     pkg.MQ        `json:"rabbitmq" yaml:"rabbitmq"`
	Oss    oss.Oss       `json:"oss" yaml:"oss"`
	Jwt    jwt.Jwt       `json:"jwt" yaml:"jwt"`
	Wechat wechat.Config `json:"wechat" yaml:"wechat"`

	WithCronJob bool `json:"-"`
	WithOutMQ   bool `json:"-"`
}

var (
	runtimeConfig  *App
	configFilePath string
)

func ReadConfig(configFilePath string) *App {
	var (
		configs *App
		err     error
	)

	// 设置配置文件
	viper.SetConfigType("yaml")
	viper.SetConfigFile(configFilePath)
	// 读取配置文件到结构体
	if err = viper.ReadInConfig(); err != nil {
		panic(err)
	}
	if err = viper.Unmarshal(&configs); err != nil {
		panic(err)
	}
	fmt.Println("读取配置文件...")
	SetConfig(configs)

	return configs
}

func SetConfig(cfg *App) {
	runtimeConfig = cfg
	pkg.SetConfig(pkg.Config{
		Env:    runtimeConfig.Env,
		Logger: runtimeConfig.Logger,
		MQ:     runtimeConfig.MQ,
		Oss:    runtimeConfig.Oss,
		Jwt:    runtimeConfig.Jwt,
		Wechat: runtimeConfig.Wechat,
	})
}

func GetConfig() *App {
	return runtimeConfig
}

func SetConfigPath(path string) {
	configFilePath = path
}

func GetConfigPath() string {
	return configFilePath
}
