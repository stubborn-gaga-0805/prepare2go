package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/stubborn-gaga-0805/prepare2go/configs/conf"
	"github.com/stubborn-gaga-0805/prepare2go/internal/mq"
	"github.com/stubborn-gaga-0805/prepare2go/internal/mq/sender"
	"github.com/stubborn-gaga-0805/prepare2go/internal/server/runtime"
	"github.com/stubborn-gaga-0805/prepare2go/pkg/helpers"
	"log"
	"os/signal"
	"syscall"
)

type runCmd struct {
	*baseCmd
}

type runFlags struct {
	appName        string
	appEnv         string
	appVersion     string
	withWs         bool
	withCronJob    bool
	withEtcdConfig bool
	withoutServer  bool
	withoutMQ      bool
}

var (
	flagAppName             = flag{"name", "n", "prepare-to-go", "Set application name"}
	flagAppEnvironment      = flag{"env", "e", "local", "Set the operating environment of the application"}
	flagAppVersion          = flag{"version", "v", "v1.0", "Set the version of the application"}
	flagWithWs              = flag{"with.ws", "", false, "Whether to start the websocket server"}
	flagWithCronJob         = flag{"with.cron", "", false, "Whether to start the crontab task"}
	flagAppConfig           = flag{"config", "c", "", "Set the path to the configuration file"}
	flagWithEtcdConfig      = flag{"config.etcd.enable", "", false, "Todo: Whether to use ETCD as the configuration center"}
	flagWithoutServerConfig = flag{"without.server", "", false, "Do not start the http server"}
	flagWithoutMQConfig     = flag{"without.mq", "", false, "Do not start the MQ server"}
)

func newRunCmd() *runCmd {
	sc := &runCmd{newBaseCmd()}
	sc.cmd = &cobra.Command{
		Use:     "run",
		Aliases: []string{"start", "running"},
		Short:   "Start web server (such as: http, grpc, websocket), and start Http server by default",
		Long:    `💡 Start your app... eg: aurora run -n myApp e test --with.cron --without.mq`,
		Run: func(cmd *cobra.Command, args []string) {
			sc.initRuntime(cmd)
			sc.initConfig()
			sc.initLogger()
			sc.initSystemMonitor()
			sc.run()
		},
	}
	addServerRuntimeFlag(sc.cmd, true)

	return sc
}

func (c *runCmd) run() {
	stop := make(chan error, 1)
	server, cleanup, err := initServer(c.ctx)
	if err != nil {
		panic(err)
	}
	defer cleanup()

	// 监听配置文件变化
	conf.SetConfigPath(c.configFilePath)
	go runtime.ConfigWatcher(c.ctx)

	// 启动Http服务
	go func() {
		if server.CrontabSwitchOn() {
			fmt.Println("Registering Crontab...")
			server.StartCrontab()
		}
		if !c.runFlags.withoutServer {
			fmt.Println("Starting HttpServer...")
			if err = server.Start(); err != nil {
				stop <- err
			}
		} else {
			fmt.Println("Without HttpServer")
		}
	}()

	// 启动WS服务
	if c.runFlags.withWs {
		ws, cleanup, err := initWebSocket(c.ctx)
		if err != nil {
			panic(err)
		}
		defer cleanup()

		go func() {
			fmt.Println("Starting WebSocket...")
			if err := ws.StartWithSocketIO(); err != nil {
				stop <- err
			}
		}()
	}

	// 初始化MQ
	var messageQueue *sender.Producer
	if !c.runFlags.withoutMQ {
		fmt.Println("Starting MQ Producer....")
		messageQueue = mq.NewMQProducer(c.ctx)
		fmt.Println("MQ Producer Running")

		go func() {
			fmt.Println("Starting MQ Consumer....")
			mq.NewMQConsumer(helpers.GetContextWithRequestId())
			fmt.Println("MQ Consumer Running")
		}()
	}

	// 系统监控
	go c.systemMonitor.StartMonitor()

	signalCtx, signalStop := signal.NotifyContext(c.ctx, syscall.SIGINT, syscall.SIGTERM)
	defer signalStop()

	select {
	case err := <-stop:
		log.Printf("start server error: %v", err.Error())
		return
	case <-signalCtx.Done():
		if server.CrontabSwitchOn() {
			server.StopCrontab()
		}
		if messageQueue != nil {
			messageQueue.RedisMQ.Stop()
		}
		fmt.Println("\nHttpServer Stopped...")
	}

	if err := server.Shutdown(c.ctx); err != nil {
		panic(err)
	}
}

// 通过命令注入运行环境参数
func addServerRuntimeFlag(cmd *cobra.Command, persistent bool) {
	getFlags(cmd, persistent).StringP(flagAppName.name, flagAppName.shortName, flagAppName.defaultValue.(string), flagAppName.usage)
	getFlags(cmd, persistent).StringP(flagAppEnvironment.name, flagAppEnvironment.shortName, flagAppEnvironment.defaultValue.(string), flagAppEnvironment.usage)
	getFlags(cmd, persistent).StringP(flagAppVersion.name, flagAppVersion.shortName, flagAppVersion.defaultValue.(string), flagAppVersion.usage)
	getFlags(cmd, persistent).StringP(flagAppConfig.name, flagAppConfig.shortName, flagAppConfig.defaultValue.(string), flagAppConfig.usage)
	getFlags(cmd, persistent).BoolP(flagWithEtcdConfig.name, flagWithEtcdConfig.shortName, flagWithEtcdConfig.defaultValue.(bool), flagWithEtcdConfig.usage)
	getFlags(cmd, persistent).BoolP(flagWithCronJob.name, flagWithCronJob.shortName, flagWithCronJob.defaultValue.(bool), flagWithCronJob.usage)
	getFlags(cmd, persistent).BoolP(flagWithWs.name, flagWithWs.shortName, flagWithWs.defaultValue.(bool), flagWithWs.usage)
	getFlags(cmd, persistent).BoolP(flagWithoutServerConfig.name, flagWithoutServerConfig.shortName, flagWithoutServerConfig.defaultValue.(bool), flagWithoutServerConfig.usage)
	getFlags(cmd, persistent).BoolP(flagWithoutMQConfig.name, flagWithoutMQConfig.shortName, flagWithoutMQConfig.defaultValue.(bool), flagWithoutMQConfig.usage)
}

// 从命令中获取AppName
func getAppName(cmd *cobra.Command) AppName {
	return AppName(cmd.Flag(flagAppName.name).Value.String())
}

// 从命令中获取Env
func getAppEnvironment(cmd *cobra.Command) Env {
	return Env(cmd.Flag(flagAppEnvironment.name).Value.String())
}

// 从命令中获取Version
func getAppVersion(cmd *cobra.Command) Version {
	return Version(cmd.Flag(flagAppVersion.name).Value.String())
}

// 从命令中获取ConfigPath
func getAppConfigPath(cmd *cobra.Command) ConfigFilePath {
	return ConfigFilePath(cmd.Flag(flagAppConfig.name).Value.String())
}

// 从命令中获取IsEtcdConfig
func getIsEtcdConfig(cmd *cobra.Command) bool {
	var isEtcConfig bool
	getFlags(cmd, false).BoolVar(&isEtcConfig, flagWithEtcdConfig.name, flagWithEtcdConfig.defaultValue.(bool), flagWithEtcdConfig.usage)
	return isEtcConfig
}

// 从命令中获取WithCronJob
func getWithCronJob(cmd *cobra.Command) bool {
	var withCronJob bool
	getFlags(cmd, false).BoolVar(&withCronJob, flagWithCronJob.name, flagWithCronJob.defaultValue.(bool), flagWithCronJob.usage)
	return withCronJob
}

func getWithoutServerConfig(cmd *cobra.Command) bool {
	var withoutServer bool
	getFlags(cmd, false).BoolVar(&withoutServer, flagWithoutServerConfig.name, flagWithoutServerConfig.defaultValue.(bool), flagWithoutServerConfig.usage)
	return withoutServer
}

func getWithoutMQConfig(cmd *cobra.Command) bool {
	var withoutMQ bool
	getFlags(cmd, false).BoolVar(&withoutMQ, flagWithoutMQConfig.name, flagWithoutMQConfig.defaultValue.(bool), flagWithoutMQConfig.usage)
	return withoutMQ
}
