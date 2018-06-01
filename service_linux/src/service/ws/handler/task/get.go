package task

import (
	"encoding/json"
	"log"

	"github.com/lillilli/todo_manager/service_linux/src/service/models"
)

// GetHandler - обработчик task_get сообщений
type GetHandler struct {
}

// Handle - обработать сообщение
func (t GetHandler) Handle(user *models.User, data []byte) error {
	var tasks []*models.Task

	if err := json.Unmarshal(data, &tasks); err != nil {
		log.Printf("[ERROR] get task ws handler: can`t parse params: %v", err)
		return user.SendParseError("task_get")
	}

	user.Stash.Set(tasks)
	return nil
}
