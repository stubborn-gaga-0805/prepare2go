package redis

import "time"

type Redis struct {
	Addr               string        `json:"addr" yaml:"addr"`
	Username           string        `json:"username" yaml:"username"`
	Password           string        `json:"password" yaml:"password"`
	Db                 int           `json:"db" yaml:"db"`
	ReadTimeout        time.Duration `json:"readTimeout" yaml:"readTimeout"`
	WriteTimeout       time.Duration `json:"writeTimeout" yaml:"writeTimeout"`
	DialTimeout        time.Duration `json:"dialTimeout" yaml:"dialTimeout"`
	MaxRetries         int           `json:"maxRetries" yaml:"maxRetries"`
	MinRetryBackoff    time.Duration `json:"minRetryBackoff" yaml:"minRetryBackoff"`
	MaxRetryBackoff    time.Duration `json:"maxRetryBackoff" yaml:"maxRetryBackoff"`
	PoolSize           int           `json:"poolSize" yaml:"poolSize"`
	MinIdleConns       int           `json:"minIdleConns" yaml:"minIdleConns"`
	MaxConnAge         time.Duration `json:"maxConnAge" yaml:"maxConnAge"`
	PoolTimeout        time.Duration `json:"poolTimeout" yaml:"poolTimeout"`
	IdleTimeout        time.Duration `json:"idleTimeout" yaml:"idleTimeout"`
	IdleCheckFrequency time.Duration `json:"idleCheckFrequency" yaml:"idleCheckFrequency"`
}

type MQConfig struct {
	Addr         string        `json:"addr" yaml:"addr"`
	Username     string        `json:"username" yaml:"username"`
	Password     string        `json:"password" yaml:"password"`
	Db           int           `json:"db" yaml:"db"`
	ReadTimeout  time.Duration `json:"readTimeout" yaml:"readTimeout"`
	WriteTimeout time.Duration `json:"writeTimeout" yaml:"writeTimeout"`
	DialTimeout  time.Duration `json:"dialTimeout" yaml:"dialTimeout"`
	MaxRetries   int           `json:"maxRetries" yaml:"maxRetries"`
	PoolSize     int           `json:"poolSize" yaml:"poolSize"`
	PoolTimeout  time.Duration `json:"poolTimeout" yaml:"poolTimeout"`
}
