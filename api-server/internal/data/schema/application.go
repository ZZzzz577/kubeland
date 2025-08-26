package schema

import (
	"api-server/internal/data/mixin"
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

type Application struct {
	ent.Schema
}

func (Application) Fields() []ent.Field {
	return []ent.Field{
		field.Uint64("cluster_id").
			Comment("集群ID"),
		field.String("name").
			NotEmpty().
			MaxLen(64).
			Comment("名称"),
		field.String("description").
			Default("").
			MaxLen(512).
			Comment("描述"),
	}
}

func (Application) Mixin() []ent.Mixin {
	return []ent.Mixin{
		mixin.TimeMixin{},
		mixin.SoftDeleteMixin{},
	}
}

func (Application) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("cluster_id", "name", "delete_at").
			Unique(),
	}
}

func (Application) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("cluster", Cluster.Type).
			Ref("applications").
			Field("cluster_id").
			Unique().
			Required(),
	}
}
