package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// PostTagRelation holds the schema definition for the PostTagRelation entity.
// This is a junction table for the many-to-many relationship between Post and PostTag.
type PostTagRelation struct {
	ent.Schema
}

// Fields of the PostTagRelation.
func (PostTagRelation) Fields() []ent.Field {
	return []ent.Field{
		field.Int("id").
			Positive().
			Immutable(),
		field.Int("post_id").
			Comment("文章ID"),
		field.Int("post_tag_id").
			Comment("标签ID"),
		field.Time("created_at").
			Immutable().
			Default(time.Now).
			Comment("创建时间"),
	}
}

// Edges of the PostTagRelation.
func (PostTagRelation) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("post", Post.Type).
			Ref("tag_relations").
			Field("post_id").
			Unique().
			Required(),
		edge.From("post_tag", PostTag.Type).
			Ref("post_relations").
			Field("post_tag_id").
			Unique().
			Required(),
	}
}

// Indexes of the PostTagRelation.
func (PostTagRelation) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("post_id", "post_tag_id").
			Unique(),
		index.Fields("post_id"),
		index.Fields("post_tag_id"),
	}
}

