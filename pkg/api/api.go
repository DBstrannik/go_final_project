package api

import (
	"fmt"
	"net/http"
	"time"
)

// Init регистрирует все HTTP-обработчики (роуты) для API
func Init() {
	http.HandleFunc("/api/nextdate", nextDateHandler)  // Получение следующей даты по правилу повтора
	http.HandleFunc("/api/task", taskHandler)          // CRUD над одной задачей (по ID)
	http.HandleFunc("/api/tasks", tasksHandler)        // (ожидается: список всех задач)
	http.HandleFunc("/api/task/done", doneTaskHandler) // Отметка задачи как выполненной
}

// nextDateHandler — обработка GET /api/nextdate
// Используется, чтобы рассчитать следующую дату задачи с повторением
func nextDateHandler(w http.ResponseWriter, r *http.Request) {
	nowStr := r.FormValue("now")    // текущее время (опционально)
	start := r.FormValue("date")    // дата начала задачи
	repeat := r.FormValue("repeat") // правило повтора (например, daily, weekly и т.д.)

	// Если параметр now не передан — используем текущее системное время
	var now time.Time
	var err error
	if nowStr == "" {
		now = time.Now()
	} else {
		// Парсим переданное время (ожидается формат как в DateFormat, например "20060102")
		now, err = time.Parse(DateFormat, nowStr)
		if err != nil {
			http.Error(w, "invalid now date", http.StatusBadRequest)
			return
		}
	}

	// Вычисляем следующую дату по логике повторения
	result, err := NextDate(now, start, repeat)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Отправляем результат обычным текстом (например: "20250722")
	fmt.Fprint(w, result)
}
