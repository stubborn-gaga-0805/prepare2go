package cmd

import (
	"context"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/stubborn-gaga-0805/prepare2go/configs/conf"
	"github.com/stubborn-gaga-0805/prepare2go/internal/server/runtime"
	"github.com/stubborn-gaga-0805/prepare2go/pkg/consts"
	"github.com/stubborn-gaga-0805/prepare2go/pkg/helpers"
	"github.com/stubborn-gaga-0805/prepare2go/pkg/logger"
	"os"
	"strings"
)

type baseCmd struct {
	cmd            *cobra.Command
	ctx            context.Context
	id             string
	isDebug        bool
	configFilePath string
	env            Env

	runFlags      *runFlags
	jobFlags      *jobFlags
	cronFlags     *crontabFlags
	genModelFlags *genModelFlags
	systemMonitor *runtime.SystemStatus
}

func newBaseCmd() *baseCmd {
	id, _ := os.Hostname()
	return &baseCmd{
		id:       id,
		ctx:      helpers.GetContextWithRequestId(),
		runFlags: new(runFlags),
	}
}

// 初始化运行时环境
func (bc *baseCmd) initRuntime(cmd *cobra.Command) {
	var (
		err      error
		appEnv   = getAppEnvironment(cmd)
		runFlags = &runFlags{
			appName:    string(getAppName(cmd)),
			appVersion: string(getAppVersion(cmd)),
		}
	)

	bc.id, _ = os.Hostname()
	bc.env = appEnv
	if !appEnv.Check() {
		panic(fmt.Sprintf("Unsupported running environment... 【%s】", bc.env))
	}
	bc.isDebug = appEnv.IsDebug()
	bc.runFlags = runFlags
	fmt.Println(fmt.Sprintf("Current App name: %s, running environment: %s", bc.runFlags.appName, bc.env))

	// 初始化配置信息
	configPath := getAppConfigPath(cmd)
	// fmt.Printf("%s \n", configPath)
	if configPath.UserDefined() {
		bc.configFilePath = configPath.ToString()
	} else {
		bc.configFilePath = fmt.Sprintf("./configs/config.%s.yaml", bc.env)
	}
	bc.runFlags.withEtcdConfig, err = cmd.Flags().GetBool(flagWithEtcdConfig.name)
	bc.runFlags.withCronJob, err = cmd.Flags().GetBool(flagWithCronJob.name)
	bc.runFlags.withWs, err = cmd.Flags().GetBool(flagWithWs.name)
	bc.runFlags.withoutServer, err = cmd.Flags().GetBool(flagWithoutServerConfig.name)
	bc.runFlags.withoutMQ, err = cmd.Flags().GetBool(flagWithoutMQConfig.name)
	fmt.Println(fmt.Sprintf("\nconfiguration file path：%s", bc.configFilePath))
	if err != nil {
		panic(err)
	}

	return
}

func (bc *baseCmd) initJobRuntime(cmd *cobra.Command) {
	bc.id, _ = os.Hostname()
	bc.env = Env(os.Getenv(consts.OSEnvKey))
	bc.configFilePath = fmt.Sprintf("./configs/config.%s.yaml", bc.env)
	customList, err := cmd.Flags().GetBool(flagCustomList.name)
	interaction, err := cmd.Flags().GetBool(flagInteraction.name)
	if err != nil {
		panic(err)
	}
	bc.jobFlags = &jobFlags{
		jobName:     getJobName(cmd),
		jobParams:   strings.Split(getCustomParams(cmd), ","),
		customList:  customList,
		interaction: interaction,
	}
	if len(bc.jobFlags.jobName) == 0 {
		bc.runFlags.withoutMQ = true
	}
	return
}

func (bc *baseCmd) initCrontabRuntime(cmd *cobra.Command) {
	bc.id, _ = os.Hostname()
	bc.env = Env(os.Getenv(consts.OSEnvKey))
	bc.configFilePath = fmt.Sprintf("./configs/config.%s.yaml", bc.env)
	crontabList, err := cmd.Flags().GetBool(flagCrontabList.name)
	if err != nil {
		panic(err)
	}
	bc.cronFlags = &crontabFlags{
		crontabList: crontabList,
	}
	return
}

func (bc *baseCmd) initGenModelRuntime(cmd *cobra.Command) {
	bc.id, _ = os.Hostname()
	bc.env = Env(os.Getenv(consts.OSEnvKey))
	bc.configFilePath = fmt.Sprintf("./configs/config.%s.yaml", bc.env)
	bc.genModelFlags = &genModelFlags{
		flagTables:      getTables(cmd),
		flagOutputPath:  getOutputPath(cmd),
		flagPackageName: getPackageName(cmd),
		flagDBConn:      getDbConn(cmd),
	}
	return
}

// 初始化、读取配置
func (bc *baseCmd) initConfig() {
	if bc.runFlags.withEtcdConfig {
		// Todo: ETCD conf center code...
	}
	cfg := conf.ReadConfig(bc.configFilePath)

	cfg.WithCronJob = bc.runFlags.withCronJob
	cfg.WithOutMQ = bc.runFlags.withoutMQ
	conf.SetConfig(cfg)

	return
}

func (bc *baseCmd) initLogger() {
	config := conf.GetConfig()
	logger.InitLog(config.Logger, config.Env.AppEnv == consts.EnvLocal)
	fmt.Println("Initializing logging system...")
	return
}

func (bc *baseCmd) initSystemMonitor() {
	bc.systemMonitor = runtime.NewSystemStatus(bc.ctx)
	return
}

func (bc *baseCmd) getCmd() *cobra.Command {
	return bc.cmd
}

func (bc *baseCmd) addCommands(commands ...cmder) {
	for _, command := range commands {
		bc.cmd.AddCommand(command.getCmd())
	}
}
