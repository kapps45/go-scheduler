package db

import (
	"database/sql"
	"fmt"
	"go-scheduler/pkg/env"
	"os"

	_ "modernc.org/sqlite"
)

const envDbFile = "TODO_DBFILE"

// schema — SQL для создания таблицы задач.
const schema = `
CREATE TABLE scheduler (
    id      INTEGER PRIMARY KEY AUTOINCREMENT,
    date    CHAR(8)      NOT NULL DEFAULT "",
    title   VARCHAR(256) NOT NULL DEFAULT "",
    comment TEXT         NOT NULL DEFAULT "",
    repeat  VARCHAR(128)          DEFAULT ""
)`

const createIndex = `CREATE INDEX idx_scheduler_date ON scheduler (date)`

var database *sql.DB

// Init открывает базу данных и создаёт таблицу, если файл не существовал.
func Init(dbFile string) error {
	actualFile := env.CheckEnv(envDbFile, dbFile)

	_, err := os.Stat(actualFile)
	install := os.IsNotExist(err)
	if err != nil && !install {
		return fmt.Errorf("проверка файла БД: %v", err)
	}

	db, err := sql.Open("sqlite", actualFile)
	if err != nil {
		return fmt.Errorf("открытие базы данных: %v", err)
	}

	if install {
		if _, err := db.Exec(schema); err != nil {
			db.Close()
			return fmt.Errorf("создание таблицы: %v", err)
		}
		if _, err := db.Exec(createIndex); err != nil {
			db.Close()
			return fmt.Errorf("создание индекса: %v", err)
		}
	}

	database = db
	return nil
}

func Close() error { return database.Close() }
