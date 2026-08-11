package api

import (
	"go-scheduler/pkg/db"
	"net/http"
	"strings"
	"time"
)

type TasksResp struct {
	Tasks []*db.Task `json:"tasks"`
}

// isDate проверяет, является ли строка датой в формате DD.MM.YYYY.
func isDate(search string) bool {
	normalized := strings.ReplaceAll(search, ",", ".")
	_, err := time.Parse(DBdateFormat, normalized)
	return err == nil
}

// convertDate переводит дату из формата DD.MM.YYYY в YYYYMMDD.
func convertDate(dateStr string) (string, error) {
	normalized := strings.ReplaceAll(dateStr, ",", ".")
	t, err := time.Parse(DBdateFormat, normalized)
	if err != nil {
		return "", err
	}
	return t.Format(dateformat), nil
}

// TasksHandler обрабатывает GET /api/tasks[?search=...] — возвращает список задач.
func TasksHandler(w http.ResponseWriter, r *http.Request) {
	search := r.URL.Query().Get("search")
	var result []*db.Task
	var err error

	if search != "" {
		if isDate(search) {
			date, err := convertDate(search)
			if err != nil {
				writeJson(w, http.StatusBadRequest, errResponse{"неверный формат даты"})
				return
			}
			result, err = db.DateTask(date, limit)
			if err != nil {
				writeJson(w, http.StatusInternalServerError, errResponse{"ошибка поиска по дате"})
				return
			}
		} else {
			result, err = db.SearchTask(search, limit)
			if err != nil {
				writeJson(w, http.StatusInternalServerError, errResponse{"ошибка поиска"})
				return
			}
		}
	} else {
		result, err = db.Tasks(limit)
	}

	if err != nil {
		writeJson(w, http.StatusInternalServerError, errResponse{"внутренняя ошибка сервера"})
		return
	}
	if result == nil {
		result = []*db.Task{}
	}

	writeJson(w, http.StatusOK, TasksResp{Tasks: result})
}
