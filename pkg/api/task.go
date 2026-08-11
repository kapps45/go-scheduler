package api

import (
	"go-scheduler/pkg/db"
	"net/http"
	"strings"
)

// GetTaskHandler обрабатывает GET /api/task?id=N — возвращает задачу по id.
func GetTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		writeJson(w, http.StatusBadRequest, errResponse{"не указан идентификатор"})
		return
	}
	task, err := db.GetTask(id)
	if err != nil {
		writeJson(w, http.StatusNotFound, errResponse{err.Error()})
		return
	}
	writeJson(w, http.StatusOK, task)
}

// UpdateTaskHandler обрабатывает PUT /api/task — обновляет существующую задачу.
func UpdateTaskHandler(w http.ResponseWriter, r *http.Request) {
	task, err := decodeTask(r)
	if err != nil {
		writeJson(w, http.StatusBadRequest, errResponse{err.Error()})
		return
	}
	if task.ID == "" {
		writeJson(w, http.StatusBadRequest, errResponse{"не указан идентификатор"})
		return
	}
	if err := normalizeDate(task); err != nil {
		writeJson(w, http.StatusBadRequest, errResponse{err.Error()})
		return
	}
	if err := db.UpdateTask(task); err != nil {
		status := http.StatusInternalServerError
		if strings.Contains(err.Error(), "не найдена") {
			status = http.StatusNotFound
		}
		writeJson(w, status, errResponse{err.Error()})
		return
	}
	writeJson(w, http.StatusOK, map[string]any{})
}
