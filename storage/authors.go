package storage

import (
	"github.com/utilyre/todolist/database"
	"github.com/utilyre/todolist/model"
)

type AuthorsStorage struct {
	database *database.Database
}

func NewAuthorsStorage(d *database.Database) AuthorsStorage {
	return AuthorsStorage{database: d}
}

func (s AuthorsStorage) Create(author *model.Author) error {
	query := `
	INSERT INTO "authors"
	("name", "email", "password")
	VALUES (?, ?, ?);
	`

	r, err := s.database.DB.Exec(query, author.Name, author.Email, author.Password)
	if err != nil {
		return err
	}

	id, err := r.LastInsertId()
	if err != nil {
		return err
	}

	author.ID = uint64(id)
	return nil
}

func (s AuthorsStorage) GetAll() ([]model.Author, error) {
	query := `
	SELECT "id", "name", "email"
	FROM "authors";
	`

	authors := []model.Author{}
	if err := s.database.DB.Select(&authors, query); err != nil {
		return nil, err
	}

	return authors, nil
}
