package handler

import (
	"errors"

	"github.com/lillilli/todo_manager/service_linux/src/service/models"
)

// Handler - интерфейс обработчика
type Handler interface {
	Handle(user *models.User, data []byte) error
}

// DefaultHandler - стандартный обработчик сообщений
type DefaultHandler struct{}

// Handle - обработать сообщение
func (h DefaultHandler) Handle(user *models.User, data []byte) error {
	return errors.New("unknown message type")
}
