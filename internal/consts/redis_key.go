package consts

const AdminTokenRedisKey = "login:token:%s" // 后台用户登陆token

const (
	FeedsContentHashSet       = "feeds:content:md5"          // 爬虫内容md5集合
	PutFeedsCountToday        = "feeds:put:total:%s"         // 今日爬虫内容投递数
	PutFeedsSuccessCountToday = "feeds:put:success:total:%s" // 今日爬虫内容投递成功数
)

const (
	MemberFavoriteList    = "member:favorite:list:%d"     // 用户收藏列表
	MemberFavoriteMsgList = "member:favorite:msg:list:%d" // 用户收藏消息列表
)

const (
	MemberInfoMQKey = "list:member:info" // redis消息队列示例
)
