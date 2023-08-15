package ecode

const (
	ParseParamsError  = 100001
	AdminUnauthorized = 100002
	UserNotFound      = 101002

	ChatGptApiError = 102001
)

var errorMsg = map[int]string{
	DemoErrorCode:     "Demo Error Msg",
	Success:           "ok",
	ServerError:       "Server Error",
	Unauthorized:      "登陆失效，请重新登陆！",
	ParseParamsError:  "参数解析错误",
	AdminUnauthorized: "登陆失效，请重新登陆",
	UserNotFound:      "用户不存在！",

	ChatGptApiError: "智能AI请求失败，请重试！",
}

var errorReason = map[int]string{
	UserNotFound: "请检查用户ID是否正确",
}

func GetErrorMsg(code int) (msg string) {
	msg, ok := errorMsg[code]
	if !ok {
		return DefaultErrorMsg
	}

	return msg
}

func GetErrorReason(code int) (reason string) {
	reason, ok := errorReason[code]
	if !ok {
		return DefaultErrorReason
	}

	return reason
}
