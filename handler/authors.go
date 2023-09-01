package handler

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/mattn/go-sqlite3"
	"github.com/utilyre/todolist/model"
	"github.com/utilyre/todolist/storage"
	"golang.org/x/crypto/bcrypt"
)

type SignUpAuthorHandler struct {
	storage storage.AuthorsStorage
}

func SetupSignUpAuthorHandler(r *mux.Router, s storage.AuthorsStorage) {
	r.Handle("/authors/signup", SignUpAuthorHandler{storage: s}).
		Methods(http.MethodPost).
		Headers("Content-Type", "application/json")
}

var _ http.Handler = SignUpAuthorHandler{}

func (h SignUpAuthorHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	author := new(model.Author)
	if err := author.DecodeAndValidate(r.Body); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(author.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Println("WARN: SignUpAuthorHandler:", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal Server Error"))
		return
	}

	author.Password = string(hash)

	if err := h.storage.Create(author); err != nil {
		errSqlite := new(sqlite3.Error)
		if errors.As(err, errSqlite) &&
			errors.Is(errSqlite.ExtendedCode, sqlite3.ErrConstraintUnique) {
			w.WriteHeader(http.StatusConflict)
			w.Write([]byte("Conflict"))
			return
		}

		log.Println("WARN: SignUpAuthorHandler:", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal Server Error"))
		return
	}

	author.Password = ""

	resp, err := json.Marshal(author)
	if err != nil {
		log.Println("WARN: SignUpAuthorHandler:", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal Server Error"))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(resp)
}

type GetAuthorsHandler struct {
	storage storage.AuthorsStorage
}

func SetupGetAuthorsHandler(r *mux.Router, s storage.AuthorsStorage) {
	r.Handle("/authors", GetAuthorsHandler{storage: s}).
		Methods(http.MethodGet)
}

var _ http.Handler = GetAuthorsHandler{}

func (h GetAuthorsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	authors, err := h.storage.GetAll()
	if err != nil {
		log.Println("WARN: GetAuthorsHandler:", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal Server Error"))
		return
	}

	resp, err := json.Marshal(authors)
	if err != nil {
		log.Println("WARN: GetAuthorsHandler:", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal Server Error"))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(resp)
}
