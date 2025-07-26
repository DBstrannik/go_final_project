package api

import (
	"encoding/json"
	"log"
	"net/http"
	"time"
)

// Константы формата даты
const (
	dateFormat = "20060102"
	DateFormat = "20060102"
)

// AfterNow сравнивает две даты, игнорируя время
func AfterNow(date, now time.Time) bool {
	dy, dm, dd := date.Date()
	ny, nm, nd := now.Date()
	return dy > ny || (dy == ny && (dm > nm || (dm == nm && dd > nd)))
}

// writeJSON записывает JSON ответ
func writeJSON(w http.ResponseWriter, data interface{}, statusCode int) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(statusCode)
	// Исправлено: Добавлена обработка ошибки кодирования JSON
	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Printf("json encode error: %v", err)
	}
}

// writeJSONError записывает JSON с сообщением об ошибке
func writeJSONError(w http.ResponseWriter, errorMsg string, statusCode int) {
	writeJSON(w, map[string]string{"error": errorMsg}, statusCode)
}
