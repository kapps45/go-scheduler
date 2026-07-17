package db

import (
	"database/sql"
	"errors"
	"fmt"
)

// Task — строка таблицы scheduler.
type Task struct {
	ID      string `json:"id"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}

// AddTask добавляет задачу и возвращает её id.
func AddTask(task *Task) (int, error) {
	res, err := database.Exec(
		`INSERT INTO scheduler (date, title, comment, repeat) VALUES (?, ?, ?, ?)`,
		task.Date, task.Title, task.Comment, task.Repeat,
	)
	if err != nil {
		return 0, fmt.Errorf("добавление задачи: %w", err)
	}
	id, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("получение id: %w", err)
	}
	return int(id), nil
}

// GetTask возвращает задачу по id.
func GetTask(id string) (*Task, error) {
	var t Task
	err := database.QueryRow(
		`SELECT id, date, title, comment, repeat FROM scheduler WHERE id = ?`, id,
	).Scan(&t.ID, &t.Date, &t.Title, &t.Comment, &t.Repeat)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("задача не найдена: %s", id)
	}
	if err != nil {
		return nil, fmt.Errorf("чтение задачи: %w", err)
	}
	return &t, nil
}

// UpdateTask обновляет все поля задачи.
func UpdateTask(task *Task) error {
	res, err := database.Exec(
		`UPDATE scheduler SET date=?, title=?, comment=?, repeat=? WHERE id=?`,
		task.Date, task.Title, task.Comment, task.Repeat, task.ID,
	)
	if err != nil {
		return fmt.Errorf("обновление задачи: %w", err)
	}
	n, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("rows affected: %w", err)
	}
	if n == 0 {
		return fmt.Errorf("задача не найдена: %s", task.ID)
	}
	return nil
}

// DeleteTask удаляет задачу по id.
func DeleteTask(id string) error {
	res, err := database.Exec(`DELETE FROM scheduler WHERE id = ?`, id)
	if err != nil {
		return fmt.Errorf("удаление задачи: %w", err)
	}
	n, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("rows affected: %w", err)
	}
	if n == 0 {
		return fmt.Errorf("задача не найдена: %s", id)
	}
	return nil
}

// Tasks возвращает до limit задач, отсортированных по дате.
func Tasks(limit int) ([]*Task, error) {
	if limit <= 0 || limit > 50 {
		limit = 50
	}
	rows, err := database.Query(
		`SELECT id, date, title, comment, repeat FROM scheduler ORDER BY date ASC LIMIT ?`,
		limit,
	)
	if err != nil {
		return nil, fmt.Errorf("список задач: %w", err)
	}
	defer rows.Close()
	return scanRows(rows)
}

// SearchTask ищет задачи по вхождению строки в заголовок или комментарий.
func SearchTask(search string, limit int) ([]*Task, error) {
	if limit <= 0 || limit > 50 {
		limit = 50
	}
	pattern := "%" + search + "%"
	rows, err := database.Query(
		`SELECT id, date, title, comment, repeat
		 FROM scheduler
		 WHERE title LIKE ? OR comment LIKE ?
		 ORDER BY date ASC LIMIT ?`,
		pattern, pattern, limit,
	)
	if err != nil {
		return nil, fmt.Errorf("поиск задач: %w", err)
	}
	defer rows.Close()
	return scanRows(rows)
}

// DateTask возвращает задачи на указанную дату (формат YYYYMMDD).
func DateTask(date string, limit int) ([]*Task, error) {
	if limit <= 0 || limit > 50 {
		limit = 50
	}
	rows, err := database.Query(
		`SELECT id, date, title, comment, repeat FROM scheduler WHERE date = ? ORDER BY date ASC LIMIT ?`,
		date, limit,
	)
	if err != nil {
		return nil, fmt.Errorf("задачи по дате: %w", err)
	}
	defer rows.Close()
	return scanRows(rows)
}

// UpdateTaskDate обновляет только дату задачи.
func UpdateTaskDate(id string, date string) error {
	res, err := database.Exec(`UPDATE scheduler SET date = ? WHERE id = ?`, date, id)
	if err != nil {
		return fmt.Errorf("обновление даты: %w", err)
	}
	n, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("rows affected: %w", err)
	}
	if n == 0 {
		return fmt.Errorf("задача не найдена: %s", id)
	}
	return nil
}

// scanRows преобразует sql.Rows в слайс Task.
func scanRows(rows *sql.Rows) ([]*Task, error) {
	var tasks []*Task
	for rows.Next() {
		var t Task
		if err := rows.Scan(&t.ID, &t.Date, &t.Title, &t.Comment, &t.Repeat); err != nil {
			return nil, fmt.Errorf("сканирование строки: %w", err)
		}
		tasks = append(tasks, &t)
	}
	return tasks, rows.Err()
}
