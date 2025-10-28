package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"time"
)

// AuditLog holds the schema definition for the AuditLog entity.
type AuditLog struct {
	ent.Schema
}

// Fields of the AuditLog.
func (AuditLog) Fields() []ent.Field {
	return []ent.Field{
		field.Int("id").
			Positive().
			Immutable(),
		field.Int("operator_id").
			Comment("操作者ID"),
		field.String("module").
			NotEmpty().
			Comment("操作模块"),
		field.String("action_type").
			NotEmpty().
			Comment("操作类型"),
		field.Int("product_id").
			Optional().
			Comment("产品ID"),
		field.String("details").
			Optional().
			SchemaType(map[string]string{
				"mysql":    "LONGTEXT", // For MySQL
				"postgres": "TEXT",     // For PostgreSQL
				"sqlite3":  "TEXT",     // For SQLite
			}).
			Comment("操作详情(JSON格式)"),
		field.String("ip_address").
			Optional().
			Comment("操作者IP地址"),
		field.Time("created_at").
			Default(time.Now).
			Immutable(),
	}
}

// Edges of the AuditLog.
func (AuditLog) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("operator", User.Type).
			Ref("audit_logs").
			Field("operator_id").
			Unique().
			Required(),
		edge.From("product", Product.Type).
			Ref("audit_logs").
			Field("product_id").
			Unique(),
	}
}
