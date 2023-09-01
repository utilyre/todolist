package storage

import "github.com/utilyre/todolist/database"

type Todo struct {
	ID    uint64 `json:"id" db:"id"`
	Title string `json:"title" validate:"required,max=16" db:"title"`
	Body  string `json:"body" validate:"max=1024" db:"body"`
}

type TodosStorage struct {
	database *database.Database
}

func NewTodosStorage(d *database.Database) TodosStorage {
	return TodosStorage{database: d}
}

func (s TodosStorage) Create(todo *Todo) error {
	query := `
	INSERT INTO "todos"
	("title", "body")
	VALUES (?, ?);
	`

	r, err := s.database.DB.Exec(query, todo.Title, todo.Body)
	if err != nil {
		return err
	}

	id, err := r.LastInsertId()
	if err != nil {
		return err
	}

	todo.ID = uint64(id)
	return nil
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

func (s TodosStorage) Update(todo *Todo) (uint64, error) {
	query := `
	UPDATE "todos"
	SET "title" = ?, "body" = ?
	WHERE "id" = ?;
	`

	r, err := s.database.DB.Exec(query, todo.Title, todo.Body, todo.ID)
	if err != nil {
		return 0, err
	}

	affected, err := r.RowsAffected()
	if err != nil {
		return 0, err
	}

	return uint64(affected), nil
}
