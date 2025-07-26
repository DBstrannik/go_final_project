package api

import (
	"net/http"

	dbpkg "go_final_project/pkg/db"
)

// Исправлено: лимит в константе
const tasksLimit = 50

// TasksResp представляет структуру ответа с списком задач
type TasksResp struct {
	Tasks []*dbpkg.Task `json:"tasks"`
}

// tasksHandler обрабатывает запросы на получение списка задач
func tasksHandler(w http.ResponseWriter, r *http.Request, db *dbpkg.DB) {
	// Исправлено: Добавлена проверка метода HTTP
	if r.Method != http.MethodGet {
		writeJSONError(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	tasks, err := db.Tasks(tasksLimit)
	if err != nil {
		// Исправлено: Разные коды ошибок для разных ситуаций
		writeJSONError(w, "database error", http.StatusInternalServerError)
		return
	}

	writeJSON(w, TasksResp{Tasks: tasks}, http.StatusOK)
}
