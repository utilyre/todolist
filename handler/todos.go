package handler

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	"github.com/utilyre/todolist/storage"
)

type CreateTodoHandler struct {
	storage storage.TodosStorage
}

func NewCreateTodoHandler(r *mux.Router, s storage.TodosStorage) {
	r.Handle("/todos", CreateTodoHandler{storage: s}).
		Methods(http.MethodPost).
		Headers("Content-Type", "application/json")
}

var _ http.Handler = CreateTodoHandler{}

func (h CreateTodoHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	type Body struct {
		Title string `json:"title" validate:"required,max=16"`
		Body  string `json:"body" validate:"max=1024"`
	}

	body := new(Body)
	if err := json.NewDecoder(r.Body).Decode(body); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	validate := validator.New()
	if err := validate.Struct(body); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	id, err := h.storage.Create(body.Title, body.Body)
	if err != nil {
		log.Println("WARN: CreateTodoHandler:", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal Server Error"))
		return
	}

	resp, err := json.Marshal(map[string]any{
		"id":    id,
		"title": body.Title,
		"body":  body.Body,
	})
	if err != nil {
		log.Println("WARN: CreateTodoHandler:", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal Server Error"))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(resp)
}

type GetTodosHandler struct {
	storage storage.TodosStorage
}

func NewGetTodosHandler(r *mux.Router, s storage.TodosStorage) {
	r.Handle("/todos", GetTodosHandler{storage: s}).
		Methods(http.MethodGet)
}

var _ http.Handler = GetTodosHandler{}

func (h GetTodosHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	todos, err := h.storage.GetAll()
	if err != nil {
		log.Println("WARN: GetTodosHandler:", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal Server Error"))
		return
	}

	resp, err := json.Marshal(todos)
	if err != nil {
		log.Println("WARN: GetTodosHandler:", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal Server Error"))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(resp)
}
