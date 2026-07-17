package api

import (
	"encoding/json"
	"fmt"
	"go-scheduler/pkg/db"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// decodeTask читает тело запроса и декодирует JSON в структуру Task.
func decodeTask(r *http.Request) (*db.Task, error) {
	defer r.Body.Close()
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, fmt.Errorf("ошибка чтения тела запроса")
	}
	var task db.Task
	if err := json.Unmarshal(body, &task); err != nil {
		return nil, fmt.Errorf("некорректный JSON")
	}
	return &task, nil
}

// validateRepeat проверяет корректность правила повторения.
// Допустимые значения: "" (разовая), "y" (ежегодно), "d N" (каждые N дней, 1–400).
func validateRepeat(repeat string) error {
	if repeat == "" || repeat == "y" {
		return nil
	}
	if strings.HasPrefix(repeat, "d ") {
		parts := strings.SplitN(repeat, " ", 2)
		n, err := strconv.Atoi(parts[1])
		if err != nil || n < 1 || n > 400 {
			return fmt.Errorf("количество дней должно быть от 1 до 400")
		}
		return nil
	}
	return fmt.Errorf("неподдерживаемый формат правила: %s", repeat)
}

// normalizeDate проверяет поля задачи и корректирует дату, чтобы она не была в прошлом:
// пустая дата → сегодня, прошлая дата без повторения → сегодня, прошлая + правило → следующая дата.
func normalizeDate(task *db.Task) error {
	if task.Title == "" {
		return fmt.Errorf("не указан заголовок задачи")
	}
	if err := validateRepeat(task.Repeat); err != nil {
		return err
	}

	today := NormalizeDate(time.Now())

	if task.Date == "" {
		task.Date = today.Format(dateformat)
		return nil
	}

	parsed, err := time.Parse(dateformat, task.Date)
	if err != nil {
		return fmt.Errorf("неверный формат даты, ожидается YYYYMMDD")
	}
	date := NormalizeDate(parsed)

	if !date.Before(today) {
		return nil
	}

	if task.Repeat == "" {
		task.Date = today.Format(dateformat)
		return nil
	}

	next, err := NextDate(today, task.Date, task.Repeat)
	if err != nil {
		return fmt.Errorf("не удалось вычислить следующую дату: %w", err)
	}
	task.Date = next
	return nil
}
