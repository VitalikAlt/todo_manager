package handler

import (
	"github.com/lillilli/todo_manager/service_linux/src/service/ws/handler/task"
	"github.com/lillilli/todo_manager/service_linux/src/service/ws/handler/user"
)

// NewManager - создание нового мэнеджера обработчиков
func NewManager() *Manager {
	handlers := make(map[string]Handler)
	manager := &Manager{handlers: handlers}

	manager.initializeHandlers()
	return manager
}

// Manager - мэнеджер обработчиков
type Manager struct {
	handlers map[string]Handler
}

func (m *Manager) initializeHandlers() {
	m.handlers["auth"] = &user.AuthHandler{}
	m.handlers["sign_in"] = &user.SignInHandler{}

	m.handlers["task_get"] = &task.GetHandler{}
	m.handlers["task_add"] = &task.AddHandler{}
	m.handlers["task_update"] = &task.UpdateHandler{}
	m.handlers["task_delete"] = &task.DeleteHandler{}
	m.handlers["task_reorder"] = &task.ReorderHandler{}
	m.handlers["task_complete"] = &task.CompleteHandler{}
}

// GetHander - получение обработчика сообщения по его типу, если данного обработчика не существует, то возвращается обработчик по-умолчанию
func (m Manager) GetHander(msgType string) Handler {
	handler, ok := m.handlers[msgType]

	if !ok {
		return &DefaultHandler{}
	}

	return handler
}
