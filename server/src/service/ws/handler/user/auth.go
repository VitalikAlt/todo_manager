package user

import (
	"encoding/json"
	"errors"
	"log"

	"github.com/lillilli/todo_manager/server/src/db"
	"github.com/lillilli/todo_manager/server/src/db/table"
	"github.com/lillilli/todo_manager/server/src/service/ws/hub"
)

// AuthHandler - обработчик auth сообщений
type AuthHandler struct {
	DB db.DB
}

// Handle - обработать сообщение
func (a AuthHandler) Handle(client *hub.Client, data []byte) error {
	var user table.User

	if err := json.Unmarshal(data, &user); err != nil {
		log.Printf("[ERROR] auth ws handler: can`t parse params: %v", err)
		return client.SendParseError("auth")
	}

	exist, err := a.DB.Tables().Users.Check(user.ID, user.Token)
	if err != nil {
		log.Printf("[ERROR] auth ws handler: can`t check user: %v", err)
		return errors.New("can`t check user")
	}

	if !exist {
		return errors.New("user doesn`t exist")
	}

	a.accept(client, &user)
	return nil
}

func (a AuthHandler) accept(client *hub.Client, user *table.User) {
	client.AcceptAuth(user)
	client.SendJSON("auth", "ok")
}
