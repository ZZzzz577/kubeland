package data

import (
	"api-server/internal/conf"
	"api-server/internal/data/generated"
	"api-server/internal/data/generated/migrate"
	_ "api-server/internal/data/generated/runtime"
	"context"
	"entgo.io/ent/dialect"
	"entgo.io/ent/dialect/sql"
	"github.com/google/wire"
	_ "github.com/mattn/go-sqlite3"
	"github.com/rs/zerolog/log"
)

var ProviderSet = wire.NewSet(
	NewData,
)

type Data generated.Client

func NewData(config *conf.Bootstrap) (*Data, func(), error) {
	c := config.Data
	driver := c.Database.Driver
	source := c.Database.Source
	if driver == dialect.SQLite && source == "" {
		source = "file:kubeland.db?_fk=1"
	}
	drv, err := sql.Open(
		driver,
		source,
	)
	if err != nil {
		log.Error().Err(err).Msg("failed opening connection to mysql")
		return nil, nil, err
	}
	client := generated.NewClient(generated.Driver(drv))
	cleanUp := func() {
		log.Info().Msg("closing the data resources")
		_ = client.Close()
	}
	err = client.Schema.Create(
		context.Background(),
		migrate.WithForeignKeys(false),
		migrate.WithDropColumn(true),
		migrate.WithDropIndex(true),
	)
	if err != nil {
		log.Error().Err(err).Msg("failed creating schema resources")
		return nil, cleanUp, err
	}
	return (*Data)(client), cleanUp, nil
}
