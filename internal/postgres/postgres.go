package postgres

import (
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/extra/bundebug"
)

type Config struct {
	Host        string
	DBName      string
	Password    string
	Port        string
	Debug       bool
	MaxOpenConn int
	MaxIdleConn int
	User        string
	SSLMode     string
}

func (c *Config) conn() string {
	return fmt.Sprintf(
		"dbname=%s host=%s port=%s user=%s password=%s sslmode=%s",
		c.DBName,
		c.Host,
		c.Port,
		c.User,
		c.Password,
		c.SSLMode,
	)
}

func NewDB(c *Config) (*bun.DB, error) {
	config, err := pgx.ParseConfig(c.conn())
	if err != nil {
		return nil, fmt.Errorf("parse config: %w", err)
	}
	sqldb := stdlib.OpenDB(*config)
	sqldb.SetMaxIdleConns(c.MaxIdleConn)
	sqldb.SetMaxOpenConns(c.MaxOpenConn)

	db := bun.NewDB(sqldb, pgdialect.New())
	if c.Debug {
		db.AddQueryHook(bundebug.NewQueryHook(bundebug.WithVerbose(true)))
	}

	return db, nil
}
