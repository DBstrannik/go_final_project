package api

import (
	"go_final_project/pkg/db"
	"net/http"
)

// TasksResp — структура для сериализации списка задач в JSON
type TasksResp struct {
	Tasks []*db.Task `json:"tasks"` // массив указателей на задачи
}

// tasksHandler — обработчик HTTP запроса для получения списка задач (до 50 штук)
// запрашивает задачи из базы и возвращает их клиенту в формате JSON
func tasksHandler(w http.ResponseWriter, r *http.Request) {
	tasks, err := db.Tasks(50) // запрашиваем из базы максимум 50 задач
	if err != nil {
		// при ошибке возвращаем JSON с описанием ошибки
		writeJSON(w, map[string]string{
			"error": err.Error(),
		})
		return
	}

	// при успешном получении — сериализуем список задач в JSON и отправляем клиенту
	writeJSON(w, TasksResp{
		Tasks: tasks,
	})
}
