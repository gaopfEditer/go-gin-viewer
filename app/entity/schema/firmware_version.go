package schema

import (
	"time"
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// FirmwareVersion 韧件版本表
type FirmwareVersion struct {
	ent.Schema
}

// Fields 定义表字段
func (FirmwareVersion) Fields() []ent.Field {
	return []ent.Field{
		field.Int("product_id").
			Comment("产品ID"),
		field.String("version").
			NotEmpty().
			Comment("版本号"),
		field.Time("release_date").
			Default(time.Now).
			Comment("发布日期"),
		field.String("remark").
			Optional().
			Comment("备注"),
		field.Int("created_by").
			Comment("创建人ID"),
		field.Time("created_at").
			Default(time.Now).
			Comment("创建时间"),
		field.Time("updated_at").
			Default(time.Now).
			UpdateDefault(time.Now).
			Comment("更新时间"),
	}
}

// Edges 定义关联关系
func (FirmwareVersion) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("software_versions", SoftwareVersion.Type).
			Ref("firmware_versions").
			Comment("关联的软件版本"),
		edge.From("product", Product.Type).
			Ref("firmware_versions").
			Field("product_id").
			Required().
			Unique(),
		edge.To("creator", User.Type).
			Field("created_by").
			Required().
			Unique().
			Comment("创建人"),
	}
}

// Indexes 定义索引
func (FirmwareVersion) Indexes() []ent.Index {
	return []ent.Index{
		// 确保同一产品下版本号唯一
		index.Fields("product_id", "version").
			Unique(),
	}
} 