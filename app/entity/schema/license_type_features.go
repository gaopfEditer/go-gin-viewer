package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// LicenseTypeFeatures holds the schema definition for the LicenseTypeFeatures entity.
type LicenseTypeFeatures struct {
	ent.Schema
}

// Fields of the LicenseTypeFeatures.
func (LicenseTypeFeatures) Fields() []ent.Field {
	return []ent.Field{
		field.Int("id").
			Positive().
			Immutable(),
		field.Int("license_type_id").
			Comment("许可证类型ID"),
		field.Int("feature_id").
			Comment("功能ID"),
		field.Time("created_at").
			Immutable().
			Default(time.Now),
	}
}

// Edges of the LicenseTypeFeatures.
func (LicenseTypeFeatures) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("license_type", LicenseType.Type).
			Field("license_type_id").
			Unique().
			Required(),
		edge.To("feature", ProductFeature.Type).
			Field("feature_id").
			Unique().
			Required(),
	}
} 