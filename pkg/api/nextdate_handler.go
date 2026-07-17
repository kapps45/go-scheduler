package api

import (
	"fmt"
	"net/http"
	"time"
)

// NextDateHandler обрабатывает GET /api/nextdate?now=...&date=...&repeat=...
func NextDateHandler(w http.ResponseWriter, r *http.Request) {
	nowReq := r.FormValue("now")
	dstart := r.FormValue("date")
	repeat := r.FormValue("repeat")

	if dstart == "" {
		http.Error(w, "не указан параметр date", http.StatusBadRequest)
		return
	}
	if repeat == "" {
		http.Error(w, "не указан параметр repeat", http.StatusBadRequest)
		return
	}

	var now time.Time
	var err error
	if nowReq == "" {
		now = time.Now()
	} else {
		now, err = time.Parse(dateformat, nowReq)
		if err != nil {
			http.Error(w, fmt.Sprintf("неверный параметр now: %v", err), http.StatusBadRequest)
			return
		}
	}

	nextDate, err := NextDate(now, dstart, repeat)
	if err != nil {
		http.Error(w, fmt.Sprintf("не удалось вычислить дату: %v", err), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(nextDate))
}
