package dto

// MetricEvent 增删改查所使用的 DTO

type AddMetricEvent struct {
	Type      string `json:"type" binding:"required"`    // 事件类型
	TS        int64  `json:"ts" binding:"required"`      // 上报时间戳(ms)
	Page      string `json:"page"`                       // 页面
	PageID    string `json:"pageId"`                     // 页面ID
	Referrer  string `json:"referrer"`                   // 来源
	UserAgent string `json:"userAgent"`                  // UA
	Payload   string `json:"payload" binding:"required"` // 原始JSON字符串
}

type ModifyMetricEvent struct {
	ID        int    `json:"id" binding:"required"`
	Type      string `json:"type"` // 可选修改
	TS        int64  `json:"ts"`   // 可选修改
	Page      string `json:"page"`
	PageID    string `json:"pageId"`
	Referrer  string `json:"referrer"`
	UserAgent string `json:"userAgent"`
	Payload   string `json:"payload"`
}

type MetricEventQuery struct {
	PageParams
	Type    string `form:"type" json:"type"`
	Route   string `form:"route" json:"route"` // 页面路径/名称
	PageID  string `form:"pageId" json:"pageId"`
	StartTS int64  `form:"startTs" json:"startTs"` // 开始时间戳(ms)
	EndTS   int64  `form:"endTs" json:"endTs"`     // 结束时间戳(ms)
}

type MetricEventResponse struct {
	ID        int    `json:"id"`
	Type      string `json:"type"`
	TS        int64  `json:"ts"`
	Page      string `json:"page"`
	PageID    string `json:"pageId"`
	Referrer  string `json:"referrer"`
	UserAgent string `json:"userAgent"`
	Payload   string `json:"payload"`
}
