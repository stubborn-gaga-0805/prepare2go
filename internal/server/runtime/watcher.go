package runtime

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"github.com/stubborn-gaga-0805/prepare2go/configs/conf"
	"github.com/stubborn-gaga-0805/prepare2go/internal/repo"
	"github.com/stubborn-gaga-0805/prepare2go/pkg/mysql"
	"github.com/stubborn-gaga-0805/prepare2go/pkg/rabbitmq"
	pkgredis "github.com/stubborn-gaga-0805/prepare2go/pkg/redis"
	"golang.org/x/net/context"
	"os/signal"
	"syscall"
)

func ConfigWatcher(ctx context.Context) {
	signalCtx, signalStop := signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGTERM)
	defer signalStop()

	// 用viper先读取文件
	viper.SetConfigType("yaml")
	viper.SetConfigFile(conf.GetConfigPath())

	fmt.Println("开始监听配置文件变化...")
	viper.OnConfigChange(func(in fsnotify.Event) {
		// 遍历根节点, 针对不同配置的变更要发送的不同的channel处理
		for key, _ := range viper.AllSettings() {
			switch key {
			case dataConfigKey:
				dataWatcherChan <- dataConfigKey
			case mqConfigKey:
				mqWatcherChan <- mqConfigKey
			}
		}
	})
	viper.WatchConfig()

	go listener(ctx)

	// 阻塞, 直到收到Ctrl+C, 然后关闭所有监听
	for {
		select {
		case <-signalCtx.Done():
			stopListenerChan <- true
			return
		}
	}
}

// 监听变化
func listener(ctx context.Context) {
	fmt.Println("\nConfig Listener Start...")
	for {
		select {
		case key := <-dataWatcherChan:
			dataWatcher(ctx, viper.Sub(key))
		case key := <-mqWatcherChan:
			mqWatcher(ctx, viper.Sub(key))
		case <-stopListenerChan:
			fmt.Println("\nConfig Listener Stopped...")
			return
		}
	}
}

func dataWatcher(ctx context.Context, v *viper.Viper) {
	for key, _ := range v.AllSettings() {
		// 判断mysql变更
		if v.Sub(key).IsSet("driver") {
			var db mysql.DB
			if err := v.Sub("db").Unmarshal(&db); err != nil {
				fmt.Printf("Unmarshal db err: %v\n", err)
				return
			}
			if conf.GetConfig().Data.Db.NotEquals(db) {
				fmt.Println("DB配置变动...")
				cfg := conf.GetConfig()
				cfg.Data.Db = db
				conf.SetConfig(cfg)

				repo.ReloadDB(ctx)
			}
		}
		// 判断redis变更
		if v.Sub(key).IsSet("db") {
			var redis pkgredis.Redis
			if err := v.Sub("redis").Unmarshal(&redis); err != nil {
				fmt.Printf("Unmarshal redis err: %v\n", err)
				return
			}
			if redis != conf.GetConfig().Data.Redis {
				fmt.Println("redis配置变动...")
				cfg := conf.GetConfig()
				cfg.Data.Redis = redis
				conf.SetConfig(cfg)

				repo.ReloadRedis(ctx)
			}
		}
	}
	return
}

func mqWatcher(ctx context.Context, v *viper.Viper) {
	var (
		rabbitMQ rabbitmq.Config
	)
	if err := v.Sub("rabbitMQ").Unmarshal(&rabbitMQ); err != nil {
		fmt.Printf("Unmarshal rabbitMQ err: %v\n", err)
		return
	}
	// 判断
	if rabbitMQ != conf.GetConfig().MQ.RabbitMQ {
		fmt.Println("RabbitMQ配置变动")
		cfg := conf.GetConfig()
		cfg.MQ.RabbitMQ = rabbitMQ
		conf.SetConfig(cfg)
	}
	return
}
