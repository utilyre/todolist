package main

import (
	"context"
	"log"
	"net/http"

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

			asRoute(newOKHandler),
		),
		fx.Invoke(
			fx.Annotate(
				newServeMux,
				fx.ParamTags(`group:"routes"`),
			),
		),
	).Run()
}

func asRoute(f any) any {
	return fx.Annotate(
		f,
		fx.ResultTags(`group:"routes"`),
	)
}

type OKHandler struct{}

func (h *OKHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func (h *OKHandler) Pattern() string {
	return "/ok"
}

func newOKHandler() Route {
	return &OKHandler{}
}

type Route interface {
	http.Handler

	Pattern() string
}

func newServeMux(routes []Route, lc fx.Lifecycle) *http.ServeMux {
	mux := http.NewServeMux()

	for _, route := range routes {
		mux.Handle(route.Pattern(), route)
	}

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go http.ListenAndServe(":5000", mux)
			return nil
		},
	})

	return mux
}
