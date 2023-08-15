package job

import (
	"context"
	"errors"
	"fmt"
	"github.com/apcera/termtables"
	"github.com/robfig/cron/v3"
	"github.com/stubborn-gaga-0805/prepare2go/internal/service"
	"github.com/stubborn-gaga-0805/prepare2go/pkg/helpers"
	"github.com/stubborn-gaga-0805/prepare2go/pkg/logger"
	"github.com/stubborn-gaga-0805/prepare2go/pkg/wework"
	"go.uber.org/zap"
	"runtime/debug"
	"sync"
	"time"
)

var (
	job        *Job
	StatusChan chan int
)

type ScheduleHandler func()

type CustomJobHandler func(ctx context.Context, opt ...string) error

type StableJobHandler func(ctx context.Context) error

type Job struct {
	service *service.Service
	cron    *cron.Cron

	ScheduleList  []*Schedule
	CustomJobList []*CustomJob
	StableJobList []*StableJob

	scheduleMapping map[cron.EntryID]*Schedule

	runningEntryID []cron.EntryID
	channel        *Channel
}

type Schedule struct {
	name    string
	desc    string
	crontab string
	handler ScheduleHandler
}

type CustomJob struct {
	name    string
	desc    string
	handler CustomJobHandler
}

type StableJob struct {
	name         string
	desc         string
	goroutineNum int
	handler      StableJobHandler
	ctx          context.Context
	stop         chan error
}

type Channel struct {
	updateDarenOrderStatChan chan int64
}

func NewJob(ctx context.Context) *Job {
	if job == nil {
		job = &Job{
			service: service.NewService(ctx),
		}
		{
			job.RegisterCrontab()
			job.RegisterCustomJobs()
			job.RegisterStableJobs()
			job.RegisterChannel()
		}
	}
	return job
}

// Start 启动定时任务
func (j *Job) Start() *Job {
	j.cron = cron.New(cron.WithSeconds())
	if len(j.ScheduleList) == 0 {
		fmt.Println("\n当前没有已注册的crontab...")
		return j
	}
	for _, schedule := range j.ScheduleList {
		EntryID, err := j.cron.AddFunc(schedule.crontab, schedule.handler)
		if err != nil {
			zap.S().Errorf("AddFun [%s:%s] Failed. Error: %v", schedule.crontab, schedule.name, err)
			continue
		}
		j.runningEntryID = append(j.runningEntryID, EntryID)
		j.scheduleMapping[EntryID] = schedule
	}
	j.printCrontabList()

	go j.cron.Start()

	return j
}

func (j *Job) Stop() {
	j.cron.Stop()
	fmt.Println("\nCronJob Stopped...")

	return
}

// RunCustomJob 执行用户任务
func (j *Job) RunCustomJob(ctx context.Context, wg *sync.WaitGroup, name string, params []string) error {
	defer wg.Done()
	for _, job := range j.CustomJobList {
		if job.name == name {
			logger.WithRequestId(helpers.GetRequestIdFromContext(ctx))
			return job.handler(ctx, params...)
		}
	}
	return errors.New("\n任务不存在！请检查任务名称")
}

// PrintCustomJobList 打印可执行的用户任务
func (j *Job) PrintCustomJobList() {
	if len(j.CustomJobList) == 0 {
		fmt.Print("\n当前没有可执行的用户任务...！\n")
		return
	}
	fmt.Print("当前可执行的用户任务：\n\n")
	for _, job := range j.CustomJobList {
		fmt.Println(fmt.Sprintf("[%s]: %s", job.name, job.desc))
	}
	return
}

// 打印已注册的crontab列表
func (j *Job) printCrontabList() {
	entries := j.cron.Entries()
	if len(entries) == 0 {
		fmt.Print("\n当前没有注册运行的定时任务...！\n")
		return
	}
	table := termtables.CreateTable()
	table.AddTitle("当前注册运行的定时任务")
	table.AddHeaders("任务ID", "执行周期", "定时任务名称", "定时任务描述", "上一次执行时间", "下一次执行时间", "下下次执行时间")
	table.AddSeparator()
	for _, entry := range entries {
		job, ok := j.scheduleMapping[entry.ID]
		if !ok {
			continue
		}
		table.AddRow(entry.ID, job.crontab, job.name, job.desc, entry.Prev.Format("2006-01-02 15:04:05"), entry.Schedule.Next(time.Now()).Format("2006-01-02 15:04:05"), entry.Schedule.Next(entry.Schedule.Next(time.Now())).Format("2006-01-02 15:04:05"))
	}
	fmt.Println(table.Render())
	return
}

// GetChannel 返回注册的Channel
func (j *Job) GetChannel() *Channel {
	return j.channel
}

func (j *Job) CrontabStatusWatcher() {
	fmt.Println("crontab watcher running...")
	for {
		select {
		case <-StatusChan:
			j.printCrontabList()
		default:
			//fmt.Printf("StatusChan: %d\n", len(StatusChan))
		}
		//time.Sleep(time.Second * 1)
	}
}

func (j *Job) StartStableJob() *Job {
	if len(j.StableJobList) == 0 {
		fmt.Println("\n当前没有已注册的常驻后台任务...")
		return j
	}
	for _, stable := range j.StableJobList {
		var (
			ctx = helpers.GetContextWithRequestId()
			job = stable
		)
		if stable.handler == nil {
			panic(fmt.Sprintf("常驻任务[%s] 未定义handlerFunc()", stable.name))
		}
		go func() {
			defer func() {
				if err := recover(); err != nil {
					logger.Helper().Errorf("Tmc panic: %v", err)
					wework.PanicBroadcast(err, helpers.GetRequestIdFromContext(ctx), string(debug.Stack()))
				}
			}()

			fmt.Printf("[%s]常驻任务: %s(%s) 开始执行\n", time.Now().Format(time.DateTime), job.name, job.desc)
			job.stop = make(chan error, 1)
			for {
				select {
				case err := <-job.stop:
					if err != nil {
						logger.Helper().Errorf("[%s]常驻任务: %s(%s) 执行失败！[err: %v]\n", time.Now().Format(time.DateTime), job.name, job.desc, err)
						return
					}
				default:
					logger.WithRequestId(helpers.GetRequestIdFromContext(ctx))
					job.stop <- job.handler(ctx)
					logger.Helper().Infof("[%s]常驻任务: %s(%s) 完成\n", time.Now().Format(time.DateTime), job.name, job.desc)
				}
			}
		}()
	}

	return j
}

func (j *Job) StopStableJob() {
	if len(j.StableJobList) == 0 {
		return
	}
	for _, stable := range j.StableJobList {
		stable.stop <- errors.New("用户终止任务")
	}

	fmt.Println("\nStableJob Stopped...")
	return
}
