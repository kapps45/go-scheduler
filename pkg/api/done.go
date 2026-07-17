package api

import (
	"go-scheduler/pkg/db"
	"net/http"
	"strings"
	"time"
)

// DoneTaskHandler обрабатывает POST /api/task/done?id=N — отмечает задачу выполненной.
// Разовая задача удаляется, повторяющаяся переносится на следующую дату.
func DoneTaskHandler(w http.ResponseWriter, r *http.Request) {
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
	if task.Repeat == "" {
		if err := db.DeleteTask(id); err != nil {
			writeJson(w, http.StatusInternalServerError, errResponse{err.Error()})
			return
		}
	} else {
		next, err := NextDate(time.Now(), task.Date, task.Repeat)
		if err != nil {
			writeJson(w, http.StatusInternalServerError, errResponse{err.Error()})
			return
		}
		if err := db.UpdateTaskDate(id, next); err != nil {
			writeJson(w, http.StatusInternalServerError, errResponse{err.Error()})
			return
		}
	}
	writeJson(w, http.StatusOK, map[string]any{})
}

// DeleteTaskHandler обрабатывает DELETE /api/task?id=N — удаляет задачу.
func DeleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		writeJson(w, http.StatusBadRequest, errResponse{"не указан идентификатор"})
		return
	}
	if err := db.DeleteTask(id); err != nil {
		status := http.StatusInternalServerError
		if strings.Contains(err.Error(), "не найдена") {
			status = http.StatusNotFound
		}
		writeJson(w, status, errResponse{err.Error()})
		return
	}
	writeJson(w, http.StatusOK, map[string]any{})
}
