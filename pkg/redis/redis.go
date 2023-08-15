package redis

import (
	"context"
	"github.com/go-redis/redis/v8"
)

func New(ctx context.Context, conf Redis) (*redis.Client, error) {
	option := &redis.Options{
		Addr: conf.Addr,
	}
	if conf.Username != "" {
		option.Username = conf.Username
	}
	if conf.Password != "" {
		option.Password = conf.Password
	}
	if conf.Db != 0 {
		option.DB = conf.Db
	}
	if conf.MaxRetries != 0 {
		option.MaxRetries = conf.MaxRetries
	}
	if conf.MinRetryBackoff.Seconds() != 0 {
		option.MinRetryBackoff = conf.MinRetryBackoff
	}
	if conf.MaxRetryBackoff.Seconds() != 0 {
		option.MaxRetryBackoff = conf.MaxRetryBackoff
	}
	if conf.DialTimeout.Seconds() != 0 {
		option.DialTimeout = conf.DialTimeout
	}
	if conf.ReadTimeout.Seconds() != 0 {
		option.ReadTimeout = conf.ReadTimeout
	}
	if conf.WriteTimeout.Seconds() != 0 {
		option.WriteTimeout = conf.WriteTimeout
	}
	if conf.PoolSize != 0 {
		option.PoolSize = conf.PoolSize
	}
	if conf.MinIdleConns != 0 {
		option.MinIdleConns = conf.MinIdleConns
	}
	if conf.MaxConnAge.Seconds() != 0 {
		option.MaxConnAge = conf.MaxConnAge
	}
	if conf.PoolTimeout.Seconds() != 0 {
		option.PoolTimeout = conf.PoolTimeout
	}
	if conf.IdleTimeout.Seconds() != 0 {
		option.IdleTimeout = conf.IdleTimeout
	}
	if conf.IdleCheckFrequency.Seconds() != 0 {
		option.IdleCheckFrequency = conf.IdleCheckFrequency
	}

	client := redis.NewClient(option)

	if _, err := client.Ping(ctx).Result(); err != nil {
		return nil, err
	}

	return client, nil
}
