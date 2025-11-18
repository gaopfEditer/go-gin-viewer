package entity

// RouteContext 与前端 routeCtx 对齐
type RouteContext struct {
	Page      string `json:"page"`
	PageID    string `json:"pageId"`
	Referrer  string `json:"referrer"`
	UserAgent string `json:"userAgent"`
}

// 通用 Envelope 头部字段
type EnvelopeBase struct {
	Type string       `json:"type"` // web-vitals | navigation-timing | longtask | resource-top | page-stable
	TS   int64        `json:"ts"`   // 上报时间戳（ms）
	Ctx  RouteContext `json:"ctx"`
}

// ---------- Web Vitals ----------

type WebVital struct {
	// web-vitals 库常见字段
	Name           string  `json:"name"`                     // FCP/LCP/CLS/INP/TTFB
	Value          float64 `json:"value"`                    // 指标值
	Delta          float64 `json:"delta,omitempty"`          // 相比上次的增量
	ID             string  `json:"id,omitempty"`             // web-vitals 会生成的唯一ID
	Rating         string  `json:"rating,omitempty"`         // good/needs-improvement/poor
	NavigationType string  `json:"navigationType,omitempty"` // navigate/reload/back-forward/prerender
}

type WebVitalsPayload struct {
	EnvelopeBase
	Metric WebVital `json:"metric"`
}

// ---------- Navigation Timing ----------

type NavigationDetail struct {
	StartTime                float64 `json:"startTime"`
	FetchStart               float64 `json:"fetchStart"`
	RequestStart             float64 `json:"requestStart"`
	ResponseStart            float64 `json:"responseStart"`
	ResponseEnd              float64 `json:"responseEnd"`
	DOMContentLoadedEventEnd float64 `json:"domContentLoadedEventEnd"`
	LoadEventEnd             float64 `json:"loadEventEnd"`
	TransferSize             int64   `json:"transferSize"`
	EncodedBodySize          int64   `json:"encodedBodySize"`
	DecodedBodySize          int64   `json:"decodedBodySize"`
}

type NavigationTimingPayload struct {
	EnvelopeBase
	// value 使用 nav.responseEnd，便于快速聚合
	Value  float64          `json:"value"`
	Detail NavigationDetail `json:"detail"`
}

// ---------- Long Task ----------

type LongTaskPayload struct {
	EnvelopeBase
	Value     float64 `json:"value"`               // duration
	Name      string  `json:"name,omitempty"`      // 通常为 'self'
	StartTime float64 `json:"startTime,omitempty"` // 开始时间
}

// ---------- Resource Top ----------

type ResourceItem struct {
	Name          string  `json:"name"`          // 资源URL
	InitiatorType string  `json:"initiatorType"` // script, css, img, font, etc.
	Duration      float64 `json:"duration"`      // 加载耗时
	TransferSize  int64   `json:"transferSize"`  // 传输大小
}

type ResourceTopPayload struct {
	EnvelopeBase
	Items []ResourceItem `json:"items"`
}

// ---------- Page Stable 标记 ----------

type PageStablePayload struct {
	EnvelopeBase
	Extra map[string]any `json:"extra,omitempty"` // 业务自定义补充字段
}
