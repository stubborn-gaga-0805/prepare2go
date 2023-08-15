package response

import (
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/middleware/requestid"
	"github.com/stubborn-gaga-0805/prepare2go/api/ecode"
	"github.com/stubborn-gaga-0805/prepare2go/pkg/helpers"
	"github.com/stubborn-gaga-0805/prepare2go/pkg/logger"
	"log"
	"net/http"
)

const HttpStatusOk = 200
const HttpStatusServerErr = 500

type errorResponse struct {
	Code      int    `json:"code"`
	Message   string `json:"message"`
	RequestId string `json:"request_id"`
}

type successResponse struct {
	Code      int               `json:"code"`
	Message   string            `json:"message"`
	Reason    string            `json:"reason,omitempty"`
	Data      interface{}       `json:"data,omitempty"`
	Meta      map[string]string `json:"meta,omitempty"`
	RequestId string            `json:"request_id"`
}

func SuccessJson(ctx iris.Context, data ...interface{}) {
	var err error
	traceId := requestid.Get(ctx)
	switch len(data) {
	case 0:
		err = ctx.JSON(newSuccessResponse(new(struct{}), traceId))
	case 1:
		err = ctx.JSON(newSuccessResponse(data[0], traceId))
	default:
		err = ctx.JSON(newSuccessResponse(data, traceId))
	}
	if err != nil {
		panic(err.Error())
	}
	ctx.Done()
	return
}

func ErrorJson(ctx iris.Context, e error) {
	requestId := requestid.Get(ctx)
	code := helpers.StringToInt(e.Error())
	if code == 0 {
		_ = ctx.StopWithJSON(HttpStatusServerErr, newFailedResponse(e, requestId))
	} else {
		_ = ctx.StopWithJSON(HttpStatusOk, newBizFailedResponse(ThrowErr(code), requestId))
	}
	ctx.Done()
	return
}

func SuccessString(ctx iris.Context, res string) {
	ctx.StatusCode(http.StatusOK)
	if _, err := ctx.WriteString(res); err != nil {
		logger.Helper().Errorf("WriteString error, err: %v", err)
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.Done()
		return
	}
	ctx.Done()
	return
}

func ErrorString(ctx iris.Context, err error) {
	ctx.StatusCode(http.StatusBadRequest)
	if _, err := ctx.WriteString(err.Error()); err != nil {
		logger.Helper().Errorf("WriteString error, err: %v", err)
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.Done()
		return
	}
	ctx.Done()
	return
}

func Unauthorized(ctx iris.Context) {
	ctx.StatusCode(http.StatusUnauthorized)
}

func newSuccessResponse(data interface{}, requestId string) successResponse {
	return successResponse{
		Code:      ecode.Success,
		Message:   ecode.GetErrorMsg(ecode.Success),
		Data:      data,
		RequestId: requestId,
	}
}

func newBizFailedResponse(e *Exception, requestId string) successResponse {
	log.Printf("BizError err: %+v", e)
	return successResponse{
		Code:      e.Code,
		Message:   ecode.GetErrorMsg(e.Code),
		Reason:    ecode.GetErrorReason(e.Code),
		RequestId: requestId,
	}
}

func newFailedResponse(err error, requestId string) errorResponse {
	log.Printf("SystemError err: %v", err)
	msg := err.Error()
	return errorResponse{
		Code:      ecode.ServerError,
		Message:   msg,
		RequestId: requestId,
	}
}
