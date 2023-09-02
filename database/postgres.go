package database

import (
	"context"
	"os"
	"path/filepath"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"github.com/utilyre/todolist/config"
	"go.uber.org/fx"
)

type Database struct {
	DSN string
	DB  *sqlx.DB
}

func New(lc fx.Lifecycle, c config.Config) *Database {
	database := &Database{DSN: c.DBPath}

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			if err := os.MkdirAll(filepath.Dir(database.DSN), 0700); err != nil {
				return err
			}

			db, err := sqlx.Connect("sqlite3", database.DSN)
			if err != nil {
				return err
			}

			database.DB = db
			return nil
		},
		OnStop: func(ctx context.Context) error {
			return database.DB.Close()
		},
	})

	return database
}
