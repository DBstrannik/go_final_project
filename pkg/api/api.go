package api

import (
	"net/http"

	dbpkg "go_final_project/pkg/db"
)

// Init инициализирует маршруты API
func Init(db *dbpkg.DB) {
	http.HandleFunc("/api/nextdate", nextDateHandler)
	http.HandleFunc("/api/task", func(w http.ResponseWriter, r *http.Request) {
		taskHandler(w, r, db)
	})
	http.HandleFunc("/api/tasks", func(w http.ResponseWriter, r *http.Request) {
		tasksHandler(w, r, db)
	})
	http.HandleFunc("/api/task/done", func(w http.ResponseWriter, r *http.Request) {
		doneTaskHandler(w, r, db)
	})
}
