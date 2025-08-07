package schema

import (
	"api-server/internal/data/mixin"
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
)

type Cluster struct {
	ent.Schema
}

func (Cluster) Fields() []ent.Field {
	return []ent.Field{
		field.String("name").
			NotEmpty().
			MaxLen(256).
			Comment("集群名称"),
		field.String("description").
			Default("").
			MaxLen(1024).
			Comment("集群描述"),
	}
}

func (Cluster) Edges() []ent.Edge {
	return nil
}

func (Cluster) Mixin() []ent.Mixin {
	return []ent.Mixin{
		mixin.TimeMixin{},
		mixin.SoftDeleteMixin{},
	}
}
