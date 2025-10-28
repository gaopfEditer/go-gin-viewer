package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"time"
)

// User holds the schema definition for the User entity.
type User struct {
	ent.Schema
}

// Fields of the User.
func (User) Fields() []ent.Field {
	return []ent.Field{
		field.Int("id").
			Immutable(),
		field.String("email").
			NotEmpty().
			Unique(),
		field.String("password").
			NotEmpty().
			Sensitive(), // 密码字段会在日志中隐藏
		field.Bool("is_enabled").
			Default(true).
			Comment("用户是否正常启用"),
		field.Time("last_login_at").
			Optional(),
		field.Time("created_at").
			Immutable().
			Default(time.Now), // 自动设置创建时间
		field.Time("updated_at").
			Default(time.Now).
			UpdateDefault(time.Now), // 自动更新时间
	}
}

// Edges of the User.
func (User) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("products", ProductManager.Type), // 用户管理的产品
		edge.To("audit_logs", AuditLog.Type),     // 用户的操作日志
		edge.From("created_devices", Device.Type).
			Ref("creator"),
		edge.From("updated_devices", Device.Type).
			Ref("updater"),
	}
}

// Indexes of the User.
func (User) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("email").
			Unique(),
	}
}
