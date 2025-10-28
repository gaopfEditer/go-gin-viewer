package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// Device schema.
type Device struct {
	ent.Schema
}

// Fields of the Device.
func (Device) Fields() []ent.Field {
	return []ent.Field{
		field.Int("id"),
		field.String("sn").Unique().Comment("设备序列号"),
		field.Int("product_id").Comment("所属产品ID"),
		field.Int("license_type_id").Optional().Comment("许可证类型ID"),
		field.String("oem_tag").Optional().Default("").Comment("OEM厂商标记"),
		field.String("remark").Optional().Default("").Comment("备注"),
		field.Time("created_at").Comment("创建时间"),
		field.Int("created_by").Optional().Comment("创建人ID"),
		field.Time("updated_at").Comment("更新时间"),
		field.Int("updated_by").Optional().Comment("更新人ID"),
	}
}

// Edges of the Device.
func (Device) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("product", Product.Type).
			Ref("devices").
			Unique().
			Field("product_id").
			Required(),
		edge.From("license_type", LicenseType.Type).
			Ref("devices").
			Unique().
			Field("license_type_id"),
		edge.To("creator", User.Type).
			Field("created_by").
			Unique(),
		edge.To("updater", User.Type).
			Field("updated_by").
			Unique(),
	}
}

// Indexes of the Device.
func (Device) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("sn").Unique(),
		index.Fields("product_id"),
		index.Fields("license_type_id"),
	}
} 