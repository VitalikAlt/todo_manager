package task

import (
	"encoding/json"
	"log"

	"github.com/lillilli/todo_manager/service_linux/src/service/models"
)

// AddHandler - обработчик task_add сообщений
type AddHandler struct {
}

// Handle - обработать сообщение
func (t AddHandler) Handle(user *models.User, data []byte) error {
	var task *models.Task

	if err := json.Unmarshal(data, &task); err != nil {
		log.Printf("[ERROR] add task ws handler: can`t parse params: %v", err)
		return user.SendParseError("task_add")
	}

	user.Stash.Add(task)
	return nil
}
