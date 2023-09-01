package model

import (
	"encoding/json"
	"io"

	"github.com/go-playground/validator/v10"
)

type Todo struct {
	ID    uint64 `json:"id" validate:"-" db:"id"`
	Title string `json:"title" validate:"required,max=16" db:"title"`
	Body  string `json:"body" validate:"max=1024" db:"body"`
}

func (t *Todo) DecodeAndValidate(r io.Reader) error {
	if err := json.NewDecoder(r).Decode(t); err != nil {
		return err
	}

	validate := validator.New()
	if err := validate.Struct(t); err != nil {
		return err
	}

	return nil
}
