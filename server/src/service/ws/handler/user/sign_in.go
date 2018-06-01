package user

import (
	"errors"
	"log"

	"github.com/lillilli/todo_manager/server/src/db"
	"github.com/lillilli/todo_manager/server/src/service/ws/hub"
)

// SignInHandler - обработчик sign_in сообщений
type SignInHandler struct {
	DB db.DB
}

// Handle - обработать сообщение
func (s SignInHandler) Handle(client *hub.Client, data []byte) error {
	log.Printf("[INFO] sign in ws handler: Try to sign in client")

	user, err := s.DB.Tables().Users.New()
	if err != nil {
		log.Printf("[INFO] sign in ws handler: can`t create new user: %v", err)
		return errors.New("can`t create new user")
	}

	return client.SendJSON("sign_in", user)
}
