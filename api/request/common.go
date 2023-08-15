package request

// PageInfo 分页信息结构
type PageInfo struct {
	Page       int   `json:"page" form:"page"`
	PageSize   int   `json:"page_size" form:"page_size"`
	TotalCount int64 `json:"total_count" form:"total_count"`
}
