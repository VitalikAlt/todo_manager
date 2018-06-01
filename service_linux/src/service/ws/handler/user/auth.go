package user

import (
	"log"

	"github.com/lillilli/todo_manager/service_linux/src/service/models"
)

// AuthHandler - обработчик auth сообщений
type AuthHandler struct {
}

// Handle - обработать сообщение
func (a AuthHandler) Handle(user *models.User, data []byte) error {
	log.Printf("%s", data)
	// var user table.User

	// if err := json.Unmarshal(data, &user); err != nil {
	// 	log.Printf("[ERROR] auth ws handler: can`t parse params: %v", err)
	// 	return client.SendParseError("auth")
	// }

	// exist, err := a.DB.Tables().Users.Check(user.ID, user.Token)
	// if err != nil {
	// 	log.Printf("[ERROR] auth ws handler: can`t check user: %v", err)
	// 	return errors.New("can`t check user")
	// }

	// if !exist {
	// 	return errors.New("user doesn`t exist")
	// }

	// a.accept(client, &user)
	return user.GetTasks()
}
