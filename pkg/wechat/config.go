package wechat

import "github.com/stubborn-gaga-0805/prepare2go/pkg/redis"

type Config struct {
	MiniApp MiniAppConfig `json:"miniApp" yaml:"miniApp"`
}

type MiniAppConfig struct {
	AppId      string      `json:"appId" yaml:"appId"`
	AppSecret  string      `json:"appSecret" yaml:"appSecret"`
	Token      string      `json:"token" yaml:"token"`
	AESKey     string      `json:"aesKey" yaml:"aesKey"`
	HttpDebug  bool        `json:"httpDebug" yaml:"httpDebug"`
	CacheRedis redis.Redis `json:"cacheRedis" yaml:"cacheRedis"`
}
