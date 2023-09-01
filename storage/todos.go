package storage

import "github.com/utilyre/todolist/database"

type Todo struct {
	ID    uint   `json:"id" db:"id"`
	Title string `json:"title" db:"title"`
	Body  string `json:"body" db:"body"`
}

type TodosStorage struct {
	database *database.Database
}

func NewTodosStorage(d *database.Database) TodosStorage {
	return TodosStorage{database: d}
}

func (s TodosStorage) Create(title, body string) (int64, error) {
	query := `
	INSERT INTO "todos"
	("title", "body")
	VALUES (?, ?);
	`

	r, err := s.database.DB.Exec(query, title, body)
	if err != nil {
		return 0, err
	}

	return r.LastInsertId()
}

func (s TodosStorage) GetAll() ([]Todo, error) {
	query := `
	SELECT "id", "title", "body"
	FROM "todos";
	`

	todos := []Todo{}
	if err := s.database.DB.Select(&todos, query); err != nil {
		return nil, err
	}

	return todos, nil
}
