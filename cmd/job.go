package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/stubborn-gaga-0805/prepare2go/internal/job"
	"github.com/stubborn-gaga-0805/prepare2go/internal/mq"
	"github.com/stubborn-gaga-0805/prepare2go/pkg/helpers"
	"sync"
)

type jobCmd struct {
	*baseCmd
}

type jobFlags struct {
	jobName     string
	jobParams   []string
	customList  bool
	interaction bool
}

var (
	flagJobName      = flag{"name", "n", "", "The name of the customer task... "}
	flagCustomParams = flag{"params", "p", "", "Parameters to run the command, multiple parameters are separated by \",\"... "}
	flagCustomList   = flag{"list", "l", false, "View executable customer tasks"}
	flagInteraction  = flag{"interaction", "i", false, "todo: Interactive interface to perform user tasks"}
)

func newJobCmd() *jobCmd {
	jc := &jobCmd{newBaseCmd()}
	jc.cmd = &cobra.Command{
		Use:     "job",
		Aliases: []string{"jobs", "J", "Job", "JOB"},
		Short:   "Customer Task Related Commands",
		Long:    `💡 Customer Task Related Commands, eg: aurora job myJob -p "first_param,second_param,third_param"`,
		Run: func(cmd *cobra.Command, args []string) {
			jc.initJobRuntime(cmd)
			jc.initConfig()
			jc.initLogger()
			jc.initSystemMonitor()
			jc.run()
		},
	}
	addJobRuntimeFlag(jc.cmd, true)

	return jc
}

func (jc *jobCmd) run() {
	var (
		ctx       = helpers.GetContextWithRequestId()
		wg        = new(sync.WaitGroup)
		jobHandle = job.NewJob(ctx)
	)
	if jc.jobFlags.customList {
		jobHandle.PrintCustomJobList()
		return
	}
	if jc.jobFlags.interaction {
		return
	}
	// 初始化MQ
	if !jc.runFlags.withoutMQ {
		fmt.Println("Starting MQ Producer....")
		mq.NewMQProducer(ctx)
		fmt.Println("MQ Producer Running")

		fmt.Println("Starting MQ Consumer....")
		go mq.NewMQConsumer(ctx)
		fmt.Println("MQ Consumer Running")
	}

	wg.Add(1)
	go func() {
		// 执行任务MQ
		err := jobHandle.RunCustomJob(ctx, wg, jc.jobFlags.jobName, jc.jobFlags.jobParams)
		if err != nil {
			fmt.Printf("Task[%s]execution faile！Error:[%v]\n", jc.jobFlags.jobName, err.Error())
			return
		}
		fmt.Printf("Task[%s]execution success！", jc.jobFlags.jobName)
		return
	}()

	wg.Wait()
	return
}

func addJobRuntimeFlag(cmd *cobra.Command, persistent bool) {
	getFlags(cmd, persistent).StringP(flagJobName.name, flagJobName.shortName, flagJobName.defaultValue.(string), flagJobName.usage)
	getFlags(cmd, persistent).StringP(flagCustomParams.name, flagCustomParams.shortName, flagCustomParams.defaultValue.(string), flagCustomParams.usage)
	getFlags(cmd, persistent).BoolP(flagCustomList.name, flagCustomList.shortName, flagCustomList.defaultValue.(bool), flagCustomList.usage)
	getFlags(cmd, persistent).BoolP(flagInteraction.name, flagInteraction.shortName, flagInteraction.defaultValue.(bool), flagInteraction.usage)
}

func getJobName(cmd *cobra.Command) string {
	return cmd.Flag(flagJobName.name).Value.String()
}

func getCustomParams(cmd *cobra.Command) string {
	return cmd.Flag(flagCustomParams.name).Value.String()
}

func getCustomList(cmd *cobra.Command) bool {
	var withCustomList bool
	getFlags(cmd, false).BoolVar(&withCustomList, flagCustomList.name, flagCustomList.defaultValue.(bool), flagCustomList.usage)
	return withCustomList
}

func getInteraction(cmd *cobra.Command) bool {
	var withInteraction bool
	getFlags(cmd, false).BoolVar(&withInteraction, flagInteraction.name, flagInteraction.defaultValue.(bool), flagInteraction.usage)
	return withInteraction
}
