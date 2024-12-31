package migration

import (
	"embed"

	"github.com/uptrace/bun/migrate"
)

var migrations = migrate.NewMigrations()

// New migrations.
func New() *migrate.Migrations {
	return migrations
}

//go:embed *.sql
var sqlMigrations embed.FS

func init() {
	if err := migrations.DiscoverCaller(); err != nil {
		panic(err)
	}
	if err := migrations.Discover(sqlMigrations); err != nil {
		panic(err)
	}
}
