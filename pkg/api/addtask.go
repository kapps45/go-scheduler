package api

import (
	"go-scheduler/pkg/db"
	"net/http"
)

// AddTaskHandler обрабатывает POST /api/task — создаёт новую задачу.
func AddTaskHandler(w http.ResponseWriter, r *http.Request) {
	task, err := decodeTask(r)
	if err != nil {
		writeJson(w, http.StatusBadRequest, errResponse{err.Error()})
		return
	}
	if err := normalizeDate(task); err != nil {
		writeJson(w, http.StatusBadRequest, errResponse{err.Error()})
		return
	}
	id, err := db.AddTask(task)
	if err != nil {
		writeJson(w, http.StatusInternalServerError, errResponse{"ошибка базы данных"})
		return
	}
	writeJson(w, http.StatusCreated, idResponse{ID: id})
}
