package schema

import (
	"api-server/internal/data/mixin"
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

type ClusterConnection struct {
	ent.Schema
}

func (ClusterConnection) Fields() []ent.Field {
	return []ent.Field{
		field.Uint64("cluster_id").
			Comment("集群ID"),
		field.String("address").
			NotEmpty().
			MaxLen(512).
			Comment("集群地址"),
		field.Uint8("type").
			Comment("连接类型"),
		field.Text("ca").
			Default("").
			Comment("CA证书"),
		field.Text("cert").
			Default("").
			Comment("证书"),
		field.Text("key").
			Default("").
			Comment("密钥"),
		field.Text("token").
			Default("").
			Comment("token"),
	}
}

func (ClusterConnection) Mixin() []ent.Mixin {
	return []ent.Mixin{
		mixin.TimeMixin{},
		mixin.SoftDeleteMixin{},
	}
}

func (ClusterConnection) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("cluster_id"),
	}
}

func (ClusterConnection) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("cluster", Cluster.Type).
			Ref("connection").
			Field("cluster_id").
			Unique().
			Required(),
	}
}
