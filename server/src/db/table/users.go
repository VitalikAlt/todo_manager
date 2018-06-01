package table

import (
	"crypto/rand"
	"database/sql"
	"fmt"
	"time"
)

// UsersManager - интерфейс для взаимодействия с таблицей tasks
type UsersManager interface {
	New() (*User, error)
	Check(id int, token string) (exist bool, err error)
}

// User - структура таблицы users
type User struct {
	ID        int       `json:"id,omitempty"`
	Token     string    `json:"token,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty"`
}

// NewUsersManager - возвращает инстанс нового мэнеджера таблицы users
func NewUsersManager(connection *sql.DB) UsersManager {
	return &userManager{connection}
}

type userManager struct {
	connection *sql.DB
}

func (u userManager) New() (*User, error) {
	newUser := User{Token: randToken(), CreatedAt: time.Now()}

	rows, err := u.connection.Query("INSERT INTO users (token) VALUES ($1) RETURNING id;",
		newUser.Token)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	rows.Next()
	if err = rows.Scan(&newUser.ID); err != nil {
		return nil, err
	}

	return &newUser, nil
}

func (u userManager) Check(id int, token string) (bool, error) {
	rows, err := u.connection.Query("SELECT created_at FROM users WHERE id = $1 AND token = $2;",
		id, token)

	if !rows.Next() {
		return false, err
	}

	rows.Close()
	return true, err
}

func randToken() string {
	b := make([]byte, 8)
	rand.Read(b)
	return fmt.Sprintf("%x", b)
}
