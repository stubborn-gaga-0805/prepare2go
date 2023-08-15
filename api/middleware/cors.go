package middleware

import (
	"github.com/kataras/iris/v12"
	"go.uber.org/zap"
)

// Cors 过滤
func (m *Handler) Cors(ctx iris.Context) {
	ctx.Header("Access-Control-Allow-Origin", "*")
	ctx.Header("Access-Control-Allow-Headers", "*")
	ctx.Header("Access-Control-Allow-Credentials", "true")
	ctx.Header("Referrer-Policy", "no-referrer-when-downgrade")
	ctx.Header("Access-Control-Expose-Headers", "*, Authorization, X-Authorization")
	if ctx.Method() == iris.MethodOptions {
		ctx.Header("Access-Control-Allow-Methods", "*")
		ctx.Header("Access-Control-Allow-Headers", "*")
		ctx.Header("Access-Control-Max-Age", "86400")
		ctx.StatusCode(iris.StatusNoContent)
		zap.S().Info("Cors...Option")
	}

	ctx.Next()
}
