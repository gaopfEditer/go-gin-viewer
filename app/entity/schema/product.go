package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// Product holds the schema definition for the Product entity.
type Product struct {
	ent.Schema
}

// Fields of the Product.
func (Product) Fields() []ent.Field {
	return []ent.Field{
		field.Int("id").
			Positive().
			Immutable(),
		field.String("code").
			NotEmpty().
			Comment("产品代号").
			Unique(),
		field.String("product_type").
			Default("default").
			Optional().
			Comment("产品类别"),
		field.String("product_name").
			NotEmpty().
			Comment("产品名称").
			Unique(),
		field.Time("created_at").
			Immutable().
			Default(time.Now), // 自动设置创建时间
		field.Time("updated_at").
			Default(time.Now).
			UpdateDefault(time.Now), // 自动更新时间
	}
}

// Edges of the Product.
func (Product) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("managers", ProductManager.Type),   // 产品的所有管理者
		edge.To("license_types", LicenseType.Type), // 产品的许可证类型
		edge.To("features", ProductFeature.Type),   // 产品的功能列表
		edge.To("firmware_versions", FirmwareVersion.Type), // 产品的韧件版本
		edge.To("software_versions", SoftwareVersion.Type), // 产品的软件版本
		edge.To("devices", Device.Type),
		edge.To("audit_logs", AuditLog.Type), // 产品的审计日志
	}
}

// Indexes of the Product.
func (Product) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("code").
			Unique(),
	}
}
