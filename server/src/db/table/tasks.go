package table

import (
	"database/sql"
	"time"

	"github.com/lib/pq"
)

const (
	// TaskStatusPending - статус таска "ожидание"
	TaskStatusPending = "PENDING"
	// TaskStatusComplete - статус таска "успешно завершён"
	TaskStatusComplete = "SUCCESS"
	// TaskStatusFailed - статус таска "завершён с ошибкой"
	TaskStatusFailed = "FAILED"
	// TaskStatusForbidden - статус таска "запрещён"
	TaskStatusForbidden = "FORBIDDEN"
)

// TasksManager - интерфейс для взаимодействия с таблицей tasks
type TasksManager interface {
	GetAll(uid int) ([]*Task, error)
	Add(uid int, text string, priority string, dueDate *time.Time) (*Task, error)
	Update(id int, text string, priority string, dueDate *time.Time) error
	Complete(id int) error
	Reorder(ids []int) error
	Delete(id int) error
}

type tasksManager struct {
	connection *sql.DB
}

// Task - структура таблицы tasks
type Task struct {
	ID         int        `json:"id,omitempty"`
	UserID     int        `json:"user_id,omitempty"`
	Order      int        `json:"order,omitempty"`
	Text       string     `json:"text,omitempty"`
	Priority   string     `json:"priority,omitempty"`
	DueDate    *time.Time `json:"due_date,omitempty"`
	CreatedAt  *time.Time `json:"created_at,omitempty"`
	FinishedAt *time.Time `json:"finished_at,omitempty"`
}

// NewTasksManager - возвращает инстанс нового мэнеджера таблицы tasks
func NewTasksManager(connection *sql.DB) TasksManager {
	return &tasksManager{connection}
}

func (t tasksManager) GetAll(uid int) ([]*Task, error) {
	tasks := make([]*Task, 0)

	rows, err := t.connection.Query(`SELECT * FROM tasks WHERE user_id = $1 ORDER BY "order";`, uid)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		task := new(Task)

		if err = rows.Scan(&task.ID, &task.UserID, &task.Order, &task.Text, &task.Priority, &task.DueDate, &task.CreatedAt, &task.FinishedAt); err != nil {
			return nil, err
		}

		tasks = append(tasks, task)
	}

	return tasks, nil
}

func (t tasksManager) GetCount(uid int) (int, error) {
	var count int

	rows, err := t.connection.Query("SELECT COUNT(id) FROM tasks WHERE user_id = $1;", uid)
	if err != nil {
		return 0, err
	}

	defer rows.Close()
	rows.Next()

	if err = rows.Scan(&count); err != nil {
		return 0, err
	}

	return count, nil
}

func (t tasksManager) Add(uid int, text string, priority string, dueDate *time.Time) (*Task, error) {
	newTask := Task{UserID: uid, Text: text, Priority: priority, DueDate: dueDate}

	if dueDate == nil {
		dueDate = &time.Time{}
	}

	count, err := t.GetCount(uid)
	if err != nil {
		return nil, err
	}

	rows, err := t.connection.Query(`INSERT INTO tasks (user_id, "order", text, priority, due_date) VALUES ($1, $2, $3, $4, $5) RETURNING id, created_at;`,
		uid, count, text, priority, pq.NullTime{Time: *dueDate})
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	rows.Next()
	if err = rows.Scan(&newTask.ID, &newTask.CreatedAt); err != nil {
		return nil, err
	}

	return &newTask, nil
}

func (t tasksManager) Update(id int, text string, priority string, dueDate *time.Time) error {
	if dueDate == nil {
		dueDate = &time.Time{}
	}

	rows, err := t.connection.Query("UPDATE tasks SET text = $1, priority = $2, due_date = $3 WHERE id = $4;", text, priority, pq.NullTime{Time: *dueDate}, id)
	rows.Close()

	return err
}

func (t tasksManager) Complete(id int) error {
	rows, err := t.connection.Query("UPDATE tasks SET finished_at = $1 WHERE id = $2;", time.Now(), id)
	rows.Close()

	return err
}

func (t tasksManager) Reorder(ids []int) error {
	for num, id := range ids {
		rows, err := t.connection.Query(`UPDATE tasks SET "order" = $1 WHERE id = $2`, num, id)
		if err != nil {
			return err
		}

		rows.Close()
	}

	return nil
}

func (t tasksManager) Delete(id int) error {
	rows, err := t.connection.Query("DELETE FROM tasks WHERE id = $1;", id)
	rows.Close()

	return err
}
