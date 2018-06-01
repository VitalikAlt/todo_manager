package task

import (
	"encoding/json"
	"log"

	"github.com/lillilli/todo_manager/service_linux/src/service/models"
)

// ReorderHandler - обработчик task_reorder сообщений
type ReorderHandler struct {
}

// Handle - обработать сообщение
func (t ReorderHandler) Handle(user *models.User, data []byte) error {
	var ids []int

	if err := json.Unmarshal(data, &ids); err != nil {
		log.Printf("[ERROR] reorder task ws handler: can`t parse params: %v", err)
		return user.SendParseError("task_reorder")
	}

	user.Stash.Reorder(ids)
	return nil
}
