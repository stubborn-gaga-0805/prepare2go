package response

import (
	"errors"
	"github.com/stubborn-gaga-0805/prepare2go/api/ecode"
	"strconv"
)

type Exception struct {
	error

	Code    int
	Message string
	Reason  string
}

func ThrowErr(code int, msg ...string) *Exception {
	errorMsg := ecode.GetErrorMsg(code)
	errorReason := ecode.GetErrorReason(code)
	if len(msg) != 0 {
		errorMsg = msg[0]
	}
	return &Exception{
		error:   errors.New(errorMsg),
		Code:    code,
		Message: errorMsg,
		Reason:  errorReason,
	}
}

func (e *Exception) Error() string {
	return strconv.Itoa(e.Code)
}
