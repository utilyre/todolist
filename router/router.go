package router

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"go.uber.org/fx"
)

func New(lc fx.Lifecycle) *mux.Router {
	host, ok := os.LookupEnv("BE_HOST")
	if !ok {
		log.Fatalln("ERROR: router: cannot find BE_HOST environment variable")
	}

	port, ok := os.LookupEnv("BE_PORT")
	if !ok {
		log.Fatalln("ERROR: router: cannot find BE_PORT environment variable")
	}

	router := mux.NewRouter()
	server := http.Server{
		Addr: fmt.Sprintf(
			"%s:%s",
			host, port,
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
