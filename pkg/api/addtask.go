package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	dbpkg "go_final_project/pkg/db"
)

// taskHandler обрабатывает все запросы к /api/task
func taskHandler(w http.ResponseWriter, r *http.Request, db *dbpkg.DB) {
	switch r.Method {
	case http.MethodPost:
		addTaskHandler(w, r, db)
	case http.MethodGet:
		getTaskHandler(w, r, db)
	case http.MethodPut:
		updateTaskHandler(w, r, db)
	case http.MethodDelete:
		deleteTaskHandler(w, r, db)
	default:
		writeJSONError(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

// addTaskHandler обрабатывает создание новой задачи
func addTaskHandler(w http.ResponseWriter, r *http.Request, db *dbpkg.DB) {
	var task dbpkg.Task

	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		writeJSONError(w, "invalid JSON data", http.StatusBadRequest)
		return
	}

	if task.Title == "" {
		writeJSONError(w, "task title is required", http.StatusBadRequest)
		return
	}

	if err := checkDate(&task); err != nil {
		writeJSONError(w, err.Error(), http.StatusBadRequest)
		return
	}

	id, err := db.AddTask(&task)
	if err != nil {
		writeJSONError(w, "database error", http.StatusInternalServerError)
		return
	}

	writeJSON(w, map[string]any{"id": id}, http.StatusOK)
}

// getTaskHandler обрабатывает получение задачи по ID
func getTaskHandler(w http.ResponseWriter, r *http.Request, db *dbpkg.DB) {
	id := r.URL.Query().Get("id")
	if id == "" {
		writeJSONError(w, "task ID is required", http.StatusBadRequest)
		return
	}

	task, err := db.GetTask(id)
	if err != nil {
		if err.Error() == "task not found" {
			writeJSONError(w, "task not found", http.StatusNotFound)
		} else {
			writeJSONError(w, "database error", http.StatusInternalServerError)
		}
		return
	}

	writeJSON(w, task, http.StatusOK)
}

// updateTaskHandler обрабатывает обновление задачи
func updateTaskHandler(w http.ResponseWriter, r *http.Request, db *dbpkg.DB) {
	var task dbpkg.Task

	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		writeJSONError(w, "invalid JSON data", http.StatusBadRequest)
		return
	}

	if task.Title == "" {
		writeJSONError(w, "task title is required", http.StatusBadRequest)
		return
	}

	if err := checkDate(&task); err != nil {
		writeJSONError(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := db.UpdateTask(&task); err != nil {
		if err.Error() == "task not found" {
			writeJSONError(w, "task not found", http.StatusNotFound)
		} else {
			writeJSONError(w, "database error", http.StatusInternalServerError)
		}
		return
	}

	writeJSON(w, map[string]string{}, http.StatusOK)
}

// checkDate проверяет и корректирует дату задачи
func checkDate(task *dbpkg.Task) error {
	now := time.Now()

	// Обработка "today"
	if task.Date == "today" {
		task.Date = now.Format(dateFormat)
	}

	if task.Date == "" {
		task.Date = now.Format(dateFormat)
	}

	t, err := time.Parse(dateFormat, task.Date)
	if err != nil {
		return errors.New("invalid date format")
	}

	// Для повторяющихся задач
	if task.Repeat != "" {
		if t.Before(now) {
			next, err := NextDate(now, task.Date, task.Repeat)
			if err != nil {
				return err
			}
			task.Date = next
		}
	} else {
		if t.Before(now) {
			task.Date = now.Format(dateFormat)
		}
	}

	return nil
}

// doneTaskHandler обрабатывает отметку задачи как выполненной
func doneTaskHandler(w http.ResponseWriter, r *http.Request, db *dbpkg.DB) {
	// Исправлено: Добавлена проверка метода HTTP
	if r.Method != http.MethodPost {
		writeJSONError(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	id := r.URL.Query().Get("id")
	if id == "" {
		writeJSONError(w, "task ID is required", http.StatusBadRequest)
		return
	}

	task, err := db.GetTask(id)
	if err != nil {
		if err.Error() == "task not found" {
			writeJSONError(w, "task not found", http.StatusNotFound)
		} else {
			writeJSONError(w, "database error", http.StatusInternalServerError)
		}
		return
	}

	if task.Repeat == "" {
		if err := db.DeleteTask(id); err != nil {
			writeJSONError(w, "database error", http.StatusInternalServerError)
			return
		}
		writeJSON(w, map[string]string{}, http.StatusOK)
		return
	}

	now := time.Now()
	next, err := NextDate(now, task.Date, task.Repeat)
	if err != nil {
		writeJSONError(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := db.UpdateDate(next, id); err != nil {
		writeJSONError(w, "database error", http.StatusInternalServerError)
		return
	}

	writeJSON(w, map[string]string{}, http.StatusOK)
}

// deleteTaskHandler обрабатывает удаление задачи
func deleteTaskHandler(w http.ResponseWriter, r *http.Request, db *dbpkg.DB) {
	id := r.URL.Query().Get("id")
	if id == "" {
		writeJSONError(w, "task ID is required", http.StatusBadRequest)
		return
	}

	if err := db.DeleteTask(id); err != nil {
		if err.Error() == "task not found" {
			writeJSONError(w, "task not found", http.StatusNotFound)
		} else {
			writeJSONError(w, "database error", http.StatusInternalServerError)
		}
		return
	}

	writeJSON(w, map[string]string{}, http.StatusOK)
}
