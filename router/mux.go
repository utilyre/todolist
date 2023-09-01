package router

import (
	"context"
	"net/http"

	"github.com/gorilla/mux"
	"go.uber.org/fx"
)

func New(lc fx.Lifecycle) *mux.Router {
	router := mux.NewRouter()
	server := http.Server{
		Addr:    ":5000",
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
