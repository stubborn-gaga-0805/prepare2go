package wework

import (
	"fmt"
	"github.com/stubborn-gaga-0805/prepare2go/pkg/consts"
	"github.com/stubborn-gaga-0805/prepare2go/pkg/http"
	"time"
)

const (
	MessageTypeMD   = "markdown"
	MessageTypeText = "text"
)

const (
	webhookPanicWarning = "https://qyapi.weixin.qq.com/cgi-bin/webhook/send?key=3ce87ca3-2a93-4316-8128-bb47bbf65fde"
)

const (
	goPanicWarnText = `
### 服务 **【%s】** 发生 <font color="warning">Panic</font> 事件!!!
###### 环境：<font color="comment">%s</font>
###### 时间：<font color="comment">%s</font>
###### RequestId：<font color="warning">%s</font>

> Panic [Message]: %s
> Panic [Stack]: %s
>
>
> @All
`
	goErrorWarnText = `
### 服务 **【%s】** 发生 <font color="warning">Error</font> !!!
###### 环境：<font color="comment">%s</font>
###### 时间：<font color="comment">%s</font>
###### RequestId：<font color="warning">%s</font>

> Error [Message]: %s
> Error [Stack]: %s
>
>
> @All
`
)

type RobotMessage struct {
	Content string `json:"content"`
}

type sendBody struct {
	MessageType string `json:"msgtype"`
	Markdown    struct {
		Content string `json:"content"`
	} `json:"markdown"`
}

var (
	AppName string
	AppEnv  string
)

func PanicBroadcast(err any, requestId, stack string) {
	if AppEnv != consts.EnvProd {
		return
	}
	var (
		sendBody = sendBody{
			MessageType: MessageTypeMD,
			Markdown: struct {
				Content string `json:"content"`
			}(struct{ Content string }{
				Content: fmt.Sprintf(
					goPanicWarnText,
					AppName,
					AppEnv,
					time.Now().Format(time.RFC3339Nano),
					requestId,
					err,
					stack,
				),
			}),
		}
		client = http.NewHttpClient()
		res    = new(interface{})
	)
	resp, err := client.SetUrl(webhookPanicWarning).SetBody(sendBody).Send(http.PostMethod, res)
	if err != nil {
		fmt.Printf("RobotBroadcast Error. body: %+v, err: %v\n", sendBody, err)
		return
	}
	if resp.IsError() {
		fmt.Printf("RobotBroadcast Response Error. body: %+v, resp: %+v, err: %v\n", sendBody, resp, err)
		return
	}
	return
}

func ErrorBroadcast(e error, requestId, stack string) {
	if AppEnv != consts.EnvProd {
		return
	}
	var (
		sendBody = sendBody{
			MessageType: MessageTypeMD,
			Markdown: struct {
				Content string `json:"content"`
			}(struct{ Content string }{
				Content: fmt.Sprintf(
					goErrorWarnText,
					AppName,
					AppEnv,
					time.Now().Format(time.RFC3339Nano),
					requestId,
					e.Error(),
					stack,
				),
			}),
		}
		client = http.NewHttpClient()
		res    = new(interface{})
	)
	resp, err := client.SetUrl(webhookPanicWarning).SetBody(sendBody).Send(http.PostMethod, res)
	if err != nil {
		fmt.Printf("RobotBroadcast Error. body: %+v, err: %v\n", sendBody, err)
		return
	}
	if resp.IsError() {
		fmt.Printf("RobotBroadcast Response Error. body: %+v, resp: %+v, err: %v\n", sendBody, resp, err)
		return
	}
	return
}
