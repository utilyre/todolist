package storage

import "github.com/utilyre/todolist/database"

type Todo struct {
	ID    uint64 `json:"id" db:"id"`
	Title string `json:"title" db:"title"`
	Body  string `json:"body" db:"body"`
}

type TodosStorage struct {
	database *database.Database
}

func NewTodosStorage(d *database.Database) TodosStorage {
	return TodosStorage{database: d}
}

func (s TodosStorage) Create(title, body string) (uint64, error) {
	query := `
	INSERT INTO "todos"
	("title", "body")
	VALUES (?, ?);
	`

	r, err := s.database.DB.Exec(query, title, body)
	if err != nil {
		return 0, err
	}

	id, err := r.LastInsertId()
	if err != nil {
		return 0, err
	}

	return uint64(id), nil
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

func (s TodosStorage) Get(id uint64) (*Todo, error) {
	query := `
	SELECT "id", "title", "body"
	FROM "todos"
	WHERE "id" = ?;
	`

	todo := new(Todo)
	if err := s.database.DB.Get(todo, query, id); err != nil {
		return nil, err
	}

	return todo, nil
}
