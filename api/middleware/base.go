package middleware

import (
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/middleware/requestid"
	"github.com/stubborn-gaga-0805/prepare2go/configs/conf"
	"github.com/stubborn-gaga-0805/prepare2go/pkg/consts"
	"github.com/stubborn-gaga-0805/prepare2go/pkg/logger"
)

// Wrapping 包装一些基础的上下文信息
func (m *Handler) Wrapping(ctx iris.Context) {
	var (
		cfg       = conf.GetConfig()
		requestId = requestid.Get(ctx)
	)
	logger.WithRequestId(requestId)
	// context 中添加必要上下文信息
	ctx.Values().Set(consts.ContextRequestIdKey, requestId)
	ctx.Values().Set(consts.ContextAdminTokenKey, ctx.GetHeader(consts.ReqHeaderTokenKey))
	ctx.Values().Set(consts.ContextMiniTokenKey, ctx.GetHeader(consts.ReqHeaderAppTokenKey))
	ctx.Values().Set(consts.ContextRunTimeEnvKey, cfg.Env.AppEnv)
	// 全局日志
	logger.Helper().Infof(
		"[%s: %s] Headers %+v, [Form] %+v, [Post] %+v, [Body] %+v, [MultipartForm] %+v",
		ctx.Request().Method,
		ctx.Request().RequestURI,
		ctx.Request().Header,
		ctx.Request().Form,
		ctx.Request().PostForm,
		ctx.Request().Body,
		ctx.Request().MultipartForm,
	)
	ctx.Next()
}
