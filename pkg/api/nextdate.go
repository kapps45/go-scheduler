package api

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

// NextDate вычисляет следующую дату выполнения задачи после now по правилу repeat.
// Поддерживаемые правила: "y" — ежегодно, "d N" — каждые N дней (1–400).
func NextDate(now time.Time, dstart string, repeat string) (string, error) {
	start, err := time.Parse(dateformat, dstart)
	if err != nil {
		return "", errors.New("неверный формат начальной даты")
	}
	if repeat == "" {
		return "", errors.New("правило повторения не указано")
	}
	switch {
	case repeat == "y":
		return yearlyNext(now, start), nil
	case strings.HasPrefix(repeat, "d "):
		return dailyNext(now, start, repeat)
	default:
		return "", fmt.Errorf("неподдерживаемый формат правила: %s", repeat)
	}
}

// yearlyNext сдвигает дату на год вперёд до тех пор, пока она не окажется позже now.
// 29 февраля в невисокосный год переносится на 1 марта.
func yearlyNext(now time.Time, start time.Time) string {
	isLeap := func(y int) bool {
		return y%4 == 0 && (y%100 != 0 || y%400 == 0)
	}
	month, day := start.Month(), start.Day()
	year := start.Year() + 1

	candidate := func(y int) time.Time {
		m, d := month, day
		if m == time.February && d == 29 && !isLeap(y) {
			m, d = time.March, 1
		}
		return time.Date(y, m, d, 0, 0, 0, 0, time.UTC)
	}

	next := candidate(year)
	for !next.After(now) {
		year++
		next = candidate(year)
	}
	return next.Format(dateformat)
}

// dailyNext сдвигает дату на N дней вперёд до тех пор, пока она не окажется позже now.
func dailyNext(now time.Time, start time.Time, repeat string) (string, error) {
	parts := strings.SplitN(repeat, " ", 2)
	if len(parts) != 2 {
		return "", errors.New("неверный формат правила для дней")
	}
	n, err := strconv.Atoi(parts[1])
	if err != nil {
		return "", fmt.Errorf("неверное количество дней: %w", err)
	}
	if n < 1 || n > 400 {
		return "", fmt.Errorf("количество дней должно быть от 1 до 400, получено %d", n)
	}

	next := start
	for {
		next = next.AddDate(0, 0, n)
		if next.After(now) {
			break
		}
	}
	return next.Format(dateformat), nil
}
