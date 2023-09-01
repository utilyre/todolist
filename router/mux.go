package router

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"go.uber.org/fx"
)

func New(lc fx.Lifecycle) *mux.Router {
	router := mux.NewRouter()
	server := http.Server{
		Addr: fmt.Sprintf(
			"%s:%s",
			os.Getenv("BE_HOST"), os.Getenv("BE_PORT"),
		),
		Handler: router,
	}

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go server.ListenAndServe()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			return server.Shutdown(ctx)
		},
	})

	return router
}
