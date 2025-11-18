package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// MetricEvent 存储前端性能指标上报的统一事件
type MetricEvent struct {
	ent.Schema
}

// Fields of the MetricEvent.
func (MetricEvent) Fields() []ent.Field {
	return []ent.Field{
		field.Int("id").
			Positive().
			Immutable(),
		field.String("type").
			NotEmpty().
			Comment("事件类型：web-vitals|navigation-timing|longtask|resource-top|page-stable"),
		field.Int64("ts").
			Comment("上报时间戳(ms)"),
		field.String("page").
			Default("").
			Comment("页面路由/名称"),
		field.String("page_id").
			Default("").
			Comment("页面唯一ID"),
		field.String("referrer").
			Default("").
			Comment("来源页面"),
		field.String("user_agent").
			Default("").
			Comment("UserAgent"),
		field.Text("payload").
			Comment("原始负载(JSON字符串)"),
		field.Time("created_at").
			Immutable().
			Default(time.Now),
		field.Time("updated_at").
			Default(time.Now).
			UpdateDefault(time.Now),
	}
}

// Indexes of the MetricEvent.
func (MetricEvent) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("type"),
		index.Fields("ts"),
		index.Fields("page"),
		index.Fields("page_id"),
		index.Fields("created_at"),
	}
}

