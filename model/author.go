package model

import (
	"encoding/json"
	"io"

	"github.com/go-playground/validator/v10"
)

type Author struct {
	ID       uint64 `json:"id" validate:"-" db:"id"`
	Name     string `json:"name" validate:"required,min=3,max=255" db:"name"`
	Email    string `json:"email" validate:"required,email,max=255" db:"email"`
	Password string `json:"password" validate:"required,min=8,max=255" db:"password"`
}

func (a *Author) DecodeAndValidate(r io.Reader) error {
	if err := json.NewDecoder(r).Decode(a); err != nil {
		return err
	}

	validate := validator.New()
	if err := validate.Struct(a); err != nil {
		return err
	}

	return nil
}
