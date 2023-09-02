package auth

import (
	"github.com/gorilla/sessions"
	"github.com/utilyre/todolist/config"
)

func New(c config.Config) *sessions.CookieStore {
	return sessions.NewCookieStore(c.BESecret)
}
