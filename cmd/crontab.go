package cmd

import (
	"github.com/spf13/cobra"
	"github.com/stubborn-gaga-0805/prepare2go/internal/job"
	"github.com/stubborn-gaga-0805/prepare2go/pkg/helpers"
	"github.com/stubborn-gaga-0805/prepare2go/pkg/logger"
)

type cronCmd struct {
	*baseCmd
}

type crontabFlags struct {
	crontabList bool
}

var (
	flagCrontabList = flag{"list", "l", false, "List running crontab tasks"}
)

func newCronCmd() *cronCmd {
	c := &cronCmd{newBaseCmd()}
	c.cmd = &cobra.Command{
		Use:     "cron",
		Aliases: []string{"crontab", "Cron", "Crontab", "CRON"},
		Short:   "Crontab task related commands",
		Long:    "",
		Run: func(cmd *cobra.Command, args []string) {
			c.initCrontabRuntime(cmd)
			c.initConfig()
			c.initLogger()
			c.initSystemMonitor()
			c.run()
		},
	}
	addCrontabRuntimeFlag(c.cmd, true)

	return c
}

func (jc *cronCmd) run() {
	var (
		ctx = helpers.GetContextWithRequestId()
		j   = job.NewJob(ctx)
	)
	logger.WithRequestId(helpers.GetRequestIdFromContext(ctx))
	if jc.cronFlags.crontabList {
		job.StatusChan <- 1
		return
	}
	j.Start()
	return
}

func addCrontabRuntimeFlag(cmd *cobra.Command, persistent bool) {
	getFlags(cmd, persistent).BoolP(flagCrontabList.name, flagCrontabList.shortName, flagCrontabList.defaultValue.(bool), flagCrontabList.usage)
}

func getCrontabList(cmd *cobra.Command) bool {
	var withCrontabList bool
	getFlags(cmd, false).BoolVar(&withCrontabList, flagCrontabList.name, flagCrontabList.defaultValue.(bool), flagCrontabList.usage)
	return withCrontabList
}
