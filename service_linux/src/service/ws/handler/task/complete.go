package task

import (
	"encoding/json"
	"log"

	"github.com/lillilli/todo_manager/service_linux/src/service/models"
)

// CompleteHandler - обработчик task_complete сообщений
type CompleteHandler struct {
}

// Handle - обработать сообщение
func (t CompleteHandler) Handle(user *models.User, data []byte) error {
	var task *models.Task

	if err := json.Unmarshal(data, &task); err != nil {
		log.Printf("[ERROR] complete task ws handler: can`t parse params: %v", err)
		return user.SendParseError("task_complete")
	}

	user.Stash.Delete(task.ID)
	return nil
}
