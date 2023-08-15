package logger

import (
	"github.com/sirupsen/logrus"
	"github.com/stubborn-gaga-0805/prepare2go/pkg/consts"
	"gorm.io/gorm"
	"os"
)

const (
	callBackAfterName = "after:logger"
)

// 结构体
type GormTrace struct {
}

// 实例化
func NewGormTrace() *GormTrace {
	return &GormTrace{}
}

// 实现方法一
func (plugin *GormTrace) Name() string {
	return "tracePlugin"
}

// 实现方法二
func (plugin *GormTrace) Initialize(db *gorm.DB) (err error) {
	// 增删改查均触发after函数
	_ = db.Callback().Create().After("gorm:after_create").Register(callBackAfterName, after)
	_ = db.Callback().Query().After("gorm:after_query").Register(callBackAfterName, after)
	_ = db.Callback().Delete().After("gorm:after_delete").Register(callBackAfterName, after)
	_ = db.Callback().Update().After("gorm:after_update").Register(callBackAfterName, after)
	_ = db.Callback().Row().After("gorm:row").Register(callBackAfterName, after)
	_ = db.Callback().Raw().After("gorm:raw").Register(callBackAfterName, after)
	return
}

func after(db *gorm.DB) {
	ctx := db.Statement.Context // 取ctx中的trace_id
	sql := db.Dialector.Explain(db.Statement.SQL.String(), db.Statement.Vars...)
	// 输出sql语句日志以及trace_id
	//logger.WithContext(ctx).Info(sql)
	l := &logrus.Logger{
		Out: os.Stdout,
		Formatter: &logrus.JSONFormatter{
			PrettyPrint: true,
		},
		Hooks: make(logrus.LevelHooks),
		Level: logrus.InfoLevel,
	}
	l.WithField("request_id", ctx.Value(consts.ContextRequestIdKey)).Info(sql)
	return
}
