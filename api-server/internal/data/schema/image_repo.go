package schema

import (
	"api-server/internal/data/mixin"
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"regexp"
)

type ImageRepo struct {
	ent.Schema
}

func (ImageRepo) Fields() []ent.Field {
	return []ent.Field{
		field.String("name").
			NotEmpty().
			MaxLen(64).
			Match(regexp.MustCompile("^[a-zA-Z][a-zA-Z0-9-]*[a-zA-Z0-9]$")).
			Comment("镜像仓库名称"),
		field.String("description").
			Default("").
			MaxLen(1024).
			Comment("镜像仓库描述"),
		field.String("url").
			MaxLen(512).
			Comment("镜像仓库地址"),
		field.String("username").
			Default("").
			MaxLen(256).
			Comment("镜像仓库用户名"),
		field.String("password").
			Default("").
			MaxLen(256).
			Comment("镜像仓库密码"),
	}
}

func (ImageRepo) Mixin() []ent.Mixin {
	return []ent.Mixin{
		mixin.TimeMixin{},
		mixin.SoftDeleteMixin{},
	}
}

func (ImageRepo) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("name").
			Unique(),
	}
}
