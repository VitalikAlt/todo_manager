package models

import (
	"log"
	"time"

	"github.com/gorilla/websocket"
)

// Message - протокол общения по ws
type Message struct {
	Type string      `json:"type"`
	Data interface{} `json:"data"`
}

// Task - структура таблицы tasks
type Task struct {
	ID         int        `json:"id,omitempty"`
	UserID     int        `json:"user_id,omitempty"`
	Order      int        `json:"order,omitempty"`
	Text       string     `json:"text,omitempty"`
	Priority   string     `json:"priority,omitempty"`
	DueDate    *time.Time `json:"due_date,omitempty"`
	CreatedAt  *time.Time `json:"created_at,omitempty"`
	FinishedAt *time.Time `json:"finished_at,omitempty"`
}

// User - структура таблицы users
type User struct {
	ID        int             `json:"id,omitempty"`
	Token     string          `json:"token,omitempty"`
	Conn      *websocket.Conn `json:"-"`
	Stash     *Stash          `json:"-"`
	CreatedAt time.Time       `json:"created_at,omitempty"`
}

// GetTasks - отправить сообщение о получении данных
func (u User) GetTasks() error {
	return u.SendJSON("task_get", "")
}

// GetUserID - получение id клиента
func (u User) GetUserID() int {
	return u.ID
}

// SendAuthError - отсылает клиенту ошибку аутентификации
func (u User) SendAuthError() error {
	return nil
	// return u.SendJSON("auth", "auth required")
}

// SendParseError - отсылает клиенту ошибку парсинга параметров сообщения
func (u User) SendParseError(msgType string) error {
	return nil
	// return u.SendJSON(msgType, "can`t parse message params")
}

// SendJSON - отправка json сообщения клиенту
func (u User) SendJSON(msgType string, v interface{}) error {
	if err := u.Conn.WriteJSON(Message{msgType, v}); err != nil {
		log.Printf("[ERROR] client: can`t send json to client: %v", err)
	}

	return nil
}
