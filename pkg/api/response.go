package api

import (
	"encoding/json"
	"net/http"
	"time"
)

const dateformat = "20060102"
const DBdateFormat = "02.01.2006"
const limit = 50

// errResponse — ответ при ошибке запроса.
type errResponse struct {
	Error string `json:"error"`
}

// idResponse — ответ после успешного создания задачи.
type idResponse struct {
	ID int `json:"id,omitempty"`
}

// writeJson сериализует data в JSON и отправляет ответ с указанным кодом.
func writeJson(w http.ResponseWriter, code int, data any) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(code)
	b, _ := json.Marshal(data)
	w.Write(b)
}

// NormalizeDate обнуляет время у даты (часы, минуты, секунды), оставляя только дату.
func NormalizeDate(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.UTC)
}
