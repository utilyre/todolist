package handler

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/utilyre/todolist/model"
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
	todo := new(model.Todo)
	if err := todo.DecodeAndValidate(r.Body); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	if err := h.storage.Create(todo); err != nil {
		log.Println("WARN: CreateTodoHandler:", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal Server Error"))
		return
	}

	resp, err := json.Marshal(todo)
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
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
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

type UpdateTodoHandler struct {
	storage storage.TodosStorage
}

func SetupUpdateTodoHandler(r *mux.Router, s storage.TodosStorage) {
	r.Handle("/todos/{id:[0-9]+}", UpdateTodoHandler{storage: s}).
		Methods(http.MethodPut).
		Headers("Content-Type", "application/json")
}

var _ http.Handler = UpdateTodoHandler{}

func (h UpdateTodoHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	todo := new(model.Todo)
	if err := todo.DecodeAndValidate(r.Body); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	vars := mux.Vars(r)
	id, err := strconv.ParseUint(vars["id"], 10, 0)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	todo.ID = id

	affected, err := h.storage.Update(todo)
	if err != nil {
		log.Println("WARN: UpdateTodoHandler:", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal Server Error"))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("%d rows affected", affected)))
}

type DeleteTodoHandler struct {
	storage storage.TodosStorage
}

func SetupDeleteTodoHandler(r *mux.Router, s storage.TodosStorage) {
	r.Handle("/todos/{id:[0-9]+}", DeleteTodoHandler{storage: s}).
		Methods(http.MethodDelete)
}

var _ http.Handler = DeleteTodoHandler{}

func (h DeleteTodoHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseUint(vars["id"], 10, 0)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	affected, err := h.storage.Delete(id)
	if err != nil {
		log.Println("WARN: DeleteTodoHandler:", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal Server Error"))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("%d rows affected", affected)))
}
