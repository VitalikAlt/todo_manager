package db

import (
	"database/sql"
	"fmt"
	"path"
	"runtime"

	// либа для работы с pg
	_ "github.com/lib/pq"
	"github.com/lillilli/todo_manager/server/src/config"
	"github.com/lillilli/todo_manager/server/src/db/table"

	"github.com/mattes/migrate"
	"github.com/mattes/migrate/database/postgres"

	// либа, необходимая для получения миграций из файловой системы
	_ "github.com/mattes/migrate/source/file"
)

// DB - интерфейс по взаимодействию с базой данных
type DB interface {
	Connect() error
	Tables() Tables
	Migrate(steps int) error
	Downgrade() error
	Disconnect() error
}

// db - инстанс sql базы данных
type db struct {
	Type             string
	URL              string
	maxRunnersOnHost int
	connection       *sql.DB
	tables           Tables
}

// Tables - таблицы базы данных
type Tables struct {
	Users table.UsersManager
	Tasks table.TasksManager
}

// NewDb - создание новой конфигурации бд
func NewDb(conf config.DbConnection) DB {
	db := &db{}

	db.Type = conf.Type
	db.URL = fmt.Sprintf("%s://%s:%s@%s:%d/%s?sslmode=disable", conf.Type, conf.User, conf.Password, conf.Host, conf.Port, conf.DbName)

	return db
}

// Connect - создание подключения к базе данных
func (db *db) Connect() error {
	connect, err := sql.Open(db.Type, db.URL)
	if err != nil {
		return err
	}

	db.connection = connect
	db.tables = Tables{
		Tasks: table.NewTasksManager(db.connection),
		Users: table.NewUsersManager(db.connection),
	}

	return nil
}

// Tables - получение всех таблиц базы данных
func (db db) Tables() Tables {
	return db.tables
}

// Migrate - провести миграции до заданной версии или до последней
//		steps	<int>	 указывает количество шагов миграции к примененеию. При значении 0 применятся все доступные миграции
func (db db) Migrate(steps int) error {
	m, err := db.prepareMigrations()
	if err != nil {
		return err
	}

	fmt.Print(steps == 0)

	if steps == 0 {
		err = m.Up()
		if err != nil {
			return err
		}
	} else {
		err = m.Steps(steps)
		if err != nil {
			return err
		}
	}

	return db.endMigrations(m)
}

// Downgrade - откатить все миграции
func (db db) Downgrade() error {
	m, err := db.prepareMigrations()
	if err != nil {
		return err
	}

	err = m.Down()
	if err != nil {
		return err
	}

	return db.endMigrations(m)
}

func (db db) prepareMigrations() (*migrate.Migrate, error) {
	driver, err := postgres.WithInstance(db.connection, &postgres.Config{})
	if err != nil {
		return nil, err
	}

	_, filename, _, _ := runtime.Caller(0)
	migrationsPath := path.Join(path.Dir(filename), "migrations")

	return migrate.NewWithDatabaseInstance(fmt.Sprintf("file://%s", migrationsPath), "postgres", driver)
}

// endMigrations - закрывает открытые файлы миграций и закрывает соединение с бд
func (db db) endMigrations(m *migrate.Migrate) error {
	sourceErr, dbErr := m.Close()
	if sourceErr != nil {
		return sourceErr
	}

	return dbErr
}

// Disconnect - закрытие соединения с бд
func (db db) Disconnect() error {
	return db.connection.Close()
}
