package handler

import (
	"errors"

	"github.com/lillilli/todo_manager/server/src/service/ws/hub"
)

// Handler - интерфейс обработчика
type Handler interface {
	Handle(client *hub.Client, data []byte) error
}

// DefaultHandler - стандартный обработчик сообщений
type DefaultHandler struct{}

// Handle - обработать сообщение
func (h DefaultHandler) Handle(client *hub.Client, data []byte) error {
	return errors.New("unknown message type")
}
