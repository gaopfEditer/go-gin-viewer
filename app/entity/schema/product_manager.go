package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// ProductManager holds the schema definition for the ProductManager entity.
type ProductManager struct {
	ent.Schema
}

func (ProductManager) Fields() []ent.Field {
	return []ent.Field{
		field.Int("id").
			Positive().
			Immutable(),
		field.Enum("role").
			Values("main", "assistant"), // 区分主管理员和协作者
		field.Int("user_id"),
		field.Int("product_id"),
		field.Enum("permissions").
			Optional().
			Values("read", "full").
			Default("read").
			Comment("权限：只读、完全，【主管理员和super_user(id:0)不受此字段限制】"),
		field.String("remark").
			Optional().
			Comment("备注信息"),
		field.Time("created_at").
			Immutable().
			Default(time.Now),
		field.Time("updated_at").
			Default(time.Now).
			UpdateDefault(time.Now), // 自动更新时间
	}
}

func (ProductManager) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("user", User.Type).
			Ref("products").
			Field("user_id"). // 显式指定外键字段
			Unique().
			Required(),
		edge.From("product", Product.Type).
			Ref("managers").
			Field("product_id"). // 显式指定外键字段
			Unique().
			Required(),
	}
}

//func (ProductManager) Indexes() []ent.Index {
//	return []ent.Index{
//		index.Fields("product_id", "role").
//			Unique().
//			Annotations(
//				entsql.IndexWhere("role = 'main'"), // 确保每个产品只有一个主管理员
//			),
//	}
//}
