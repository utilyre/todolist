package router

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/utilyre/todolist/config"
	"go.uber.org/fx"
)

func New(lc fx.Lifecycle, c config.Config) *mux.Router {
	router := mux.NewRouter()
	server := http.Server{
		Addr: fmt.Sprintf(
			"%s:%s",
			c.BEHost, c.BEPort,
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
