package job

// RegisterCrontab 注册定时任务
// 验证工具：http://www.jsons.cn/quartzcheck/
func (j *Job) RegisterCrontab() *Job {
	// 在这里注册crontab
	j.ScheduleList = []*Schedule{}
	return j
}

// RegisterCustomJobs 注册用户命令行执行的方法
func (j *Job) RegisterCustomJobs() *Job {
	// 在这里注册命令
	j.CustomJobList = []*CustomJob{
		{"RedisMqProducerDemo", "测试redisMQ生产者", RedisMqProduce},
	}
	return j
}

// RegisterStableJobs 注册常驻命令执行的方法
func (j *Job) RegisterStableJobs() *Job {
	// 在这里注册常驻命令
	j.StableJobList = []*StableJob{}
	return j
}

func (j *Job) RegisterChannel() *Job {
	j.channel = &Channel{
		updateDarenOrderStatChan: make(chan int64, 0),
	}
	return j
}
