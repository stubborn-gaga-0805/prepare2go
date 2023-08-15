package consts

const (
	EnvLocal = "local"
	EnvDev   = "dev"
	EnvTest  = "test"
	EnvPre   = "pre"
	EnvProd  = "prod"

	OSEnvKey = "RUNTIME_ENV"

	ContextRequestIdKey  = "CtxRequestId"
	ContextRunTimeEnvKey = "CtxRuntimeEnv"
	ContextUserIdKey     = "CtxUserId"
	ContextAdminTokenKey = "CtxAdminToken"
	ContextMiniTokenKey  = "CtxMiniAuthToken"

	ReqHeaderTokenKey     = "ADMIN-TOKEN"
	ReqHeaderAppTokenKey  = "APP-TOKEN"
	ReqHeaderRequestIdKey = "X-Request-Id"

	IsDeleted    = 1 // 删除
	IsNotDeleted = 2 // 未删除

	ResetAdminUserPasswordKey = "=sKh7B49=86zSKD7rEE4yss8P#HmTsfkzCsctJQ27d36zbHpvLu$2PC4QuJL5QPt" // 重制管理员密码key（开发测试用）

	NumericAlphabet = "1234567890"
	NormalAlphabet  = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"
	UserNoAlphabet  = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"
)
