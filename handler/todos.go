package handler

import (
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	"github.com/utilyre/todolist/storage"
)

type CreateTodoHandler struct {
	storage storage.TodosStorage
}

func SetupCreateTodoHandler(r *mux.Router, s storage.TodosStorage) {
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

func SetupGetTodosHandler(r *mux.Router, s storage.TodosStorage) {
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

type GetTodoHandler struct {
	storage storage.TodosStorage
}

func SetupGetTodoHandler(r *mux.Router, s storage.TodosStorage) {
	r.Handle("/todos/{id:[0-9]+}", GetTodoHandler{storage: s}).
		Methods(http.MethodGet)
}

var _ http.Handler = GetTodoHandler{}

func (h GetTodoHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseUint(vars["id"], 10, 0)
	if err != nil {
		log.Println("WARN: GetTodoHandler:", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal Server Error"))
		return
	}

	todo, err := h.storage.Get(id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("Not Found"))
			return
		}

		log.Println("WARN: GetTodoHandler:", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal Server Error"))
		return
	}

	resp, err := json.Marshal(todo)
	if err != nil {
		log.Println("WARN: GetTodoHandler:", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal Server Error"))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(resp)
}
