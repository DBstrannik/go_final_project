package api

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// nextDateHandler обрабатывает запросы для расчета следующей даты
func nextDateHandler(w http.ResponseWriter, r *http.Request) {
	// Исправлено: Добавлена проверка метода HTTP
	if r.Method != http.MethodGet {
		writeJSONError(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	nowStr := r.FormValue("now")
	start := r.FormValue("date")
	repeat := r.FormValue("repeat")

	var now time.Time
	var err error
	if nowStr == "" {
		now = time.Now()
	} else {
		now, err = time.Parse(DateFormat, nowStr)
		if err != nil {
			writeJSONError(w, "invalid now date", http.StatusBadRequest)
			return
		}
	}

	result, err := NextDate(now, start, repeat)
	if err != nil {
		writeJSONError(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Исправлено: Добавлена проверка ошибки при записи ответа
	if _, err := fmt.Fprint(w, result); err != nil {
		log.Printf("write response error: %v", err)
	}
}

// NextDate вычисляет следующую дату выполнения задачи
func NextDate(now time.Time, dstart string, repeat string) (string, error) {
	if repeat == "" {
		return "", errors.New("repeat rule is required")
	}

	startDate, err := time.Parse(DateFormat, dstart)
	if err != nil {
		return "", errors.New("invalid start date")
	}

	parts := strings.Split(repeat, " ")

	switch parts[0] {
	case "d":
		if len(parts) != 2 {
			return "", errors.New("invalid daily rule format")
		}
		days, err := strconv.Atoi(parts[1])
		if err != nil || days < 1 || days > 400 {
			return "", errors.New("invalid number of days")
		}
		for {
			startDate = startDate.AddDate(0, 0, days)
			if AfterNow(startDate, now) {
				break
			}
		}
		return startDate.Format(DateFormat), nil

	case "y":
		for {
			startDate = startDate.AddDate(1, 0, 0)
			if AfterNow(startDate, now) {
				break
			}
		}
		return startDate.Format(DateFormat), nil

	default:
		return "", errors.New("unsupported repeat rule")
	}
}
