package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// Post holds the schema definition for the Post entity.
type Post struct {
	ent.Schema
}

// Fields of the Post.
func (Post) Fields() []ent.Field {
	return []ent.Field{
		field.Int("id").
			Positive().
			Immutable(),
		field.String("title").
			NotEmpty().
			Comment("文章标题"),
		field.Text("content").
			Optional().
			Comment("文章内容"),
		field.Text("excerpt").
			Optional().
			Comment("文章摘要"),
		field.String("slug").
			Optional().
			Unique().
			Comment("文章别名，用于URL"),
		field.String("status").
			Default("draft").
			Comment("文章状态：draft(草稿)、published(已发布)、archived(已归档)"),
		field.Int("view_count").
			Default(0).
			Comment("浏览次数"),
		field.Time("published_at").
			Optional().
			Comment("发布时间"),
		field.Time("created_at").
			Immutable().
			Default(time.Now), // 自动设置创建时间
		field.Time("updated_at").
			Default(time.Now).
			UpdateDefault(time.Now), // 自动更新时间
	}
}

// Edges of the Post.
func (Post) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("category", PostCategory.Type).
			Ref("posts").
			Unique().
			Comment("文章所属类别"),
		edge.To("tag_relations", PostTagRelation.Type).
			Comment("文章标签关系"),
		edge.From("author", User.Type).
			Ref("posts").
			Unique().
			Comment("文章作者"),
	}
}

// Indexes of the Post.
func (Post) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("slug").
			Unique(),
		index.Fields("status"),
		index.Fields("published_at"),
		index.Fields("created_at"),
	}
}
