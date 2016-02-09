package api

import (
	"models"
)

func NewSession() *Session {
	model := &Session{}
	return model
}

type Session struct {
	models.Session
}
