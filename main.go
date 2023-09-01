package main

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/utilyre/todolist/database"
	"github.com/utilyre/todolist/handler"
	"github.com/utilyre/todolist/router"
	"github.com/utilyre/todolist/storage"
	"go.uber.org/fx"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatalln("ERROR: godotenv:", err)
	}

	fx.New(
		fx.Provide(
			database.New,
			router.New,
			storage.NewTodosStorage,
		),
		fx.Invoke(
			handler.SetupCreateTodoHandler,
			handler.SetupGetTodosHandler,
			handler.SetupGetTodoHandler,
			handler.SetupUpdateTodoHandler,
		),
	).Run()
}
