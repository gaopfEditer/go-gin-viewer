package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// LicenseType holds the schema definition for the LicenseType entity.
type LicenseType struct {
	ent.Schema
}

// Fields of the LicenseType.
func (LicenseType) Fields() []ent.Field {
	return []ent.Field{
		field.Int("id").
			Positive().
			Immutable(),
		field.String("type_name").
			NotEmpty().
			Comment("许可证类型名称"),
		field.String("license_type").
			NotEmpty().
			Immutable().
			Comment("许可证编码"),
		field.Int("product_id").
			Comment("所属产品ID"),
		field.Time("created_at").
			Immutable().
			Default(time.Now),
		field.Time("updated_at").
			Default(time.Now).
			UpdateDefault(time.Now),
	}
}

// Edges of the LicenseType.
func (LicenseType) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("product", Product.Type).
			Ref("license_types").
			Field("product_id").
			Unique().
			Required(),
		edge.To("features", ProductFeature.Type).
			Through("license_type_features", LicenseTypeFeatures.Type),
		edge.To("devices", Device.Type),
	}
}
