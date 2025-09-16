package schema

import (
	"api-server/internal/data/mixin"
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"regexp"
)

type GitRepo struct {
	ent.Schema
}

func (GitRepo) Fields() []ent.Field {
	return []ent.Field{
		field.String("name").
			NotEmpty().
			MaxLen(64).
			Match(regexp.MustCompile("^[a-zA-Z][a-zA-Z0-9-]*[a-zA-Z0-9]$")).
			Comment("git仓库名称"),
		field.String("description").
			Default("").
			MaxLen(1024).
			Comment("git仓库描述"),
		field.String("url").
			MaxLen(512).
			Comment("git仓库地址"),
		field.String("token").
			Default("").
			MaxLen(1024).
			Comment("git仓库token"),
	}
}

func (GitRepo) Mixin() []ent.Mixin {
	return []ent.Mixin{
		mixin.TimeMixin{},
		mixin.SoftDeleteMixin{},
	}
}

func (GitRepo) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("name").
			Unique(),
	}
}
