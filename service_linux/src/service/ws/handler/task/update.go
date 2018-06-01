package task

import (
	"encoding/json"
	"log"

	"github.com/lillilli/todo_manager/service_linux/src/service/models"
)

// UpdateHandler - обработчик task_add сообщений
type UpdateHandler struct {
}

// Handle - обработать сообщение
func (t UpdateHandler) Handle(user *models.User, data []byte) error {
	var task *models.Task

	if err := json.Unmarshal(data, &task); err != nil {
		log.Printf("[ERROR] add task ws handler: can`t parse params: %v", err)
		return user.SendParseError("task_add")
	}

	user.Stash.Update(task)
	return nil
}
