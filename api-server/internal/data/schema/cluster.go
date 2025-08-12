package schema

import (
	"api-server/internal/data/mixin"
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"regexp"
)

type Cluster struct {
	ent.Schema
}

func (Cluster) Fields() []ent.Field {
	return []ent.Field{
		field.String("name").
			NotEmpty().
			MaxLen(64).
			Match(regexp.MustCompile("^[a-zA-Z][a-zA-Z0-9-]*[a-zA-Z0-9]$")).
			Comment("集群名称"),
		field.String("description").
			Default("").
			MaxLen(1024).
			Comment("集群描述"),
	}
}

func (Cluster) Mixin() []ent.Mixin {
	return []ent.Mixin{
		mixin.TimeMixin{},
		mixin.SoftDeleteMixin{},
	}
}

func (Cluster) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("name").
			Unique(),
	}
}

func (Cluster) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("connection", ClusterConnection.Type).
			Unique(),
	}
}
