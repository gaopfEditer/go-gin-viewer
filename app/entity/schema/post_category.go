package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// PostCategory holds the schema definition for the PostCategory entity.
type PostCategory struct {
	ent.Schema
}

// Fields of the PostCategory.
func (PostCategory) Fields() []ent.Field {
	return []ent.Field{
		field.Int("id").
			Positive().
			Immutable(),
		field.String("name").
			NotEmpty().
			Comment("类别名称"),
		field.String("slug").
			Optional().
			Unique().
			Comment("类别别名，用于URL"),
		field.Text("description").
			Optional().
			Comment("类别描述"),
		field.String("color").
			Optional().
			Comment("类别颜色，用于前端显示"),
		field.Int("sort_order").
			Default(0).
			Comment("排序顺序"),
		field.Bool("is_active").
			Default(true).
			Comment("是否启用"),
		field.Time("created_at").
			Immutable().
			Default(time.Now), // 自动设置创建时间
		field.Time("updated_at").
			Default(time.Now).
			UpdateDefault(time.Now), // 自动更新时间
	}
}

// Edges of the PostCategory.
func (PostCategory) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("posts", Post.Type).
			Comment("该类别下的文章"),
	}
}

// Indexes of the PostCategory.
func (PostCategory) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("slug").
			Unique(),
		index.Fields("is_active"),
		index.Fields("sort_order"),
	}
}
