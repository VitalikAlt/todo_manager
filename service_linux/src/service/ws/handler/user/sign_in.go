package user

import (
	"encoding/json"
	"log"

	"github.com/lillilli/todo_manager/service_linux/src/service/models"
)

// SignInHandler - обработчик sign_in сообщений
type SignInHandler struct {
}

// Handle - обработать сообщение
func (s SignInHandler) Handle(user *models.User, data []byte) error {
	var userData models.User
	log.Printf("[INFO] sign in ws handler: Try to sign in client")

	if err := json.Unmarshal(data, &userData); err != nil {
		log.Printf("[ERROR] auth ws handler: can`t parse params: %v", err)
		return user.SendParseError("auth")
	}

	user.ID = userData.ID
	user.Token = userData.Token

	return user.GetTasks()
}
