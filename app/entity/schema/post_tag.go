package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// PostTag holds the schema definition for the PostTag entity.
type PostTag struct {
	ent.Schema
}

// Fields of the PostTag.
func (PostTag) Fields() []ent.Field {
	return []ent.Field{
		field.Int("id").
			Positive().
			Immutable(),
		field.String("name").
			NotEmpty().
			Comment("标签名称"),
		field.String("slug").
			Optional().
			Unique().
			Comment("标签别名，用于URL"),
		field.Text("description").
			Optional().
			Comment("标签描述"),
		field.String("color").
			Optional().
			Comment("标签颜色，用于前端显示"),
		field.Int("post_count").
			Default(0).
			Comment("使用该标签的文章数量"),
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

// Edges of the PostTag.
func (PostTag) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("post_relations", PostTagRelation.Type).
			Comment("标签文章关系"),
	}
}

// Indexes of the PostTag.
func (PostTag) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("slug").
			Unique(),
		index.Fields("is_active"),
		index.Fields("post_count"),
	}
}
