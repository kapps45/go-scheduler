package api

import "net/http"

// Init регистрирует все маршруты API.
func Init() {
	http.HandleFunc("/api/nextdate", NextDateHandler)
	http.HandleFunc("/api/tasks", TasksHandler)
	http.HandleFunc("/api/task/done", DoneTaskHandler)
	http.HandleFunc("/api/task", taskDispatcher)
}

// taskDispatcher маршрутизирует /api/task по HTTP-методу.
func taskDispatcher(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		GetTaskHandler(w, r)
	case http.MethodPost:
		AddTaskHandler(w, r)
	case http.MethodPut:
		UpdateTaskHandler(w, r)
	case http.MethodDelete:
		DeleteTaskHandler(w, r)
	default:
		http.Error(w, "метод не поддерживается", http.StatusMethodNotAllowed)
	}
}
