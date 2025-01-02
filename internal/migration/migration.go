package migration

import (
	"embed"

	"github.com/uptrace/bun/migrate"
)

// Migrations variable for running migrations.
var Migrations = migrate.NewMigrations()

// New returns the migrations created.
func New() *migrate.Migrations {
	return Migrations
}

//go:embed *.sql
var sqlMigrations embed.FS

func init() {
	if err := Migrations.DiscoverCaller(); err != nil {
		panic(err)
	}
	if err := Migrations.Discover(sqlMigrations); err != nil {
		panic(err)
	}
}
