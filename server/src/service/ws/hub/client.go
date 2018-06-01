package hub

import (
	"log"

	"github.com/gorilla/websocket"
	"github.com/lillilli/todo_manager/server/src/db/table"
)

// Message - протокол общения сервера  клиентов
type Message struct {
	Type string      `json:"type"`
	Data interface{} `json:"data"`
}

// Client is a middleman between the websocket connection and the hub.
type Client struct {
	hub *Hub

	// The websocket connection.
	conn *websocket.Conn

	// user, which connection is
	user *table.User
}

// AcceptAuth - подтверждает аутентификацию
func (c *Client) AcceptAuth(user *table.User) {
	c.user = user
	c.hub.clientsByID[c.user.ID] = append(c.hub.clientsByID[c.user.ID], c)
}

// CheckAuth - проверка аутентификации
func (c Client) CheckAuth() bool {
	return c.user != nil
}

// GetUserID - получение id клиента
func (c Client) GetUserID() int {
	return c.user.ID
}

// SendAuthError - отсылает клиенту ошибку аутентификации
func (c Client) SendAuthError() error {
	return c.SendJSON("auth", "auth required")
}

// SendParseError - отсылает клиенту ошибку парсинга параметров сообщения
func (c Client) SendParseError(msgType string) error {
	return c.SendJSON(msgType, "can`t parse message params")
}

// SendJSONbyUID - отправка сообщения всем клиентам с user_id = user_id текущего клиента
func (c Client) SendJSONbyUID(msgType string, v interface{}) error {
	clients, ok := c.hub.clientsByID[c.user.ID]
	if !ok {
		return nil
	}

	for _, client := range clients {
		client.SendJSON(msgType, v)
	}

	return nil
}

// SendJSON - отправка json сообщения клиенту
func (c Client) SendJSON(msgType string, v interface{}) error {
	if err := c.conn.WriteJSON(Message{msgType, v}); err != nil {
		log.Printf("[ERROR] client: can`t send json to client: %v", err)
	}
	return nil
}
