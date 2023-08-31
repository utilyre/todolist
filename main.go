package main

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/utilyre/todolist/config"
	"github.com/utilyre/todolist/database"
	"go.uber.org/fx"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatalln("ERROR: godotenv:", err)
	}

	fx.New(
		fx.Provide(
			config.NewDatabaseConfig,
			database.New,
		),
		fx.Invoke(func(database.PostgresDatabase) {}),
	).Run()
}
