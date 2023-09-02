package main

import (
	"github.com/utilyre/todolist/auth"
	"github.com/utilyre/todolist/config"
	"github.com/utilyre/todolist/database"
	"github.com/utilyre/todolist/handler"
	"github.com/utilyre/todolist/router"
	"github.com/utilyre/todolist/storage"
	"go.uber.org/fx"
)

func main() {
	fx.New(
		fx.Provide(
			config.New,
			auth.New,
			database.New,
			router.New,

			storage.NewAuthorsStorage,
			storage.NewTodosStorage,
		),
		fx.Invoke(
			handler.SetupSignUpAuthorHandler,
			handler.SetupGetAuthorsHandler,

			handler.SetupCreateTodoHandler,
			handler.SetupGetTodosHandler,
			handler.SetupGetTodoHandler,
			handler.SetupUpdateTodoHandler,
			handler.SetupDeleteTodoHandler,
		),
	).Run()
}
