package handler

import (
	"github.com/lillilli/todo_manager/server/src/db"
	"github.com/lillilli/todo_manager/server/src/service/ws/handler/task"
	"github.com/lillilli/todo_manager/server/src/service/ws/handler/user"
)

// NewManager - создание нового мэнеджера обработчиков
func NewManager(db db.DB) *Manager {
	handlers := make(map[string]Handler)
	manager := &Manager{db: db, handlers: handlers}

	manager.initializeHandlers()
	return manager
}

// Manager - мэнеджер обработчиков
type Manager struct {
	db       db.DB
	handlers map[string]Handler
}

func (m *Manager) initializeHandlers() {
	m.handlers["auth"] = &user.AuthHandler{DB: m.db}
	m.handlers["sign_in"] = &user.SignInHandler{DB: m.db}

	m.handlers["task_get"] = &task.GetHandler{DB: m.db}
	m.handlers["task_add"] = &task.AddHandler{DB: m.db}
	m.handlers["task_update"] = &task.UpdateHandler{DB: m.db}
	m.handlers["task_delete"] = &task.DeleteHandler{DB: m.db}
	m.handlers["task_reorder"] = &task.ReorderHandler{DB: m.db}
	m.handlers["task_complete"] = &task.CompleteHandler{DB: m.db}
}

// GetHander - получение обработчика сообщения по его типу, если данного обработчика не существует, то возвращается обработчик по-умолчанию
func (m Manager) GetHander(msgType string) Handler {
	handler, ok := m.handlers[msgType]

	if !ok {
		return &DefaultHandler{}
	}

	return handler
}
