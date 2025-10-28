package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// ProductFeature holds the schema definition for the ProductFeature entity.
type ProductFeature struct {
	ent.Schema
}

// Fields of the ProductFeature.
func (ProductFeature) Fields() []ent.Field {
	return []ent.Field{
		field.Int("id").
			Positive().
			Immutable(),
		field.String("feature_name").
			NotEmpty().
			Comment("功能名称"),
		field.String("feature_code").
			NotEmpty().
			Immutable().
			Comment("功能编码"),
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

// Edges of the ProductFeature.
func (ProductFeature) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("product", Product.Type).
			Ref("features").
			Field("product_id").
			Unique().
			Required(),
		edge.From("license_types", LicenseType.Type).
			Ref("features").
			Through("license_type_features", LicenseTypeFeatures.Type),
		edge.From("software_versions", SoftwareVersion.Type).
			Ref("features"),
	}
} 