package handler

import (
	"net/http"

	"github.com/gorilla/mux"
)

type CreateTodoHandler struct{}

func NewCreateTodoHandler(r *mux.Router) {
	r.Handle(
		"/todos/ok",
		CreateTodoHandler{},
	).Methods(http.MethodGet)
}

var _ http.Handler = CreateTodoHandler{}

func (h CreateTodoHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}
