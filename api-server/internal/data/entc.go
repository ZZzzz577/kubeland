//go:build ignore
// +build ignore

package main

import (
	"entgo.io/ent/entc"
	"entgo.io/ent/entc/gen"
	"entgo.io/ent/schema/field"
	"github.com/rs/zerolog/log"
)

func main() {
	if err := entc.Generate("./schema", &gen.Config{
		Target:  "generated",
		Package: "api-server/internal/data/generated",
		IDType:  &field.TypeInfo{Type: field.TypeUint64},
		Features: []gen.Feature{
			gen.FeatureIntercept,
			gen.FeatureSnapshot,
		},
	}); err != nil {
		log.Fatal().Err(err).Msg("running ent codegen")
	}
}
