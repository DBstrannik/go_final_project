package server

import (
	"go_final_project/pkg/api"
	"go_final_project/pkg/db"
	"log"
	"net/http"
	"os"
)

// Server представляет HTTP сервер приложения
type Server struct {
	logger *log.Logger
	port   string
	db     *db.DB
}

// NewServer создает новый экземпляр сервера
func NewServer(logger *log.Logger, db *db.DB) *Server {
	port := os.Getenv("TODO_PORT")
	if port == "" {
		port = "7540" // Порт по умолчанию
	}

	return &Server{
		logger: logger,
		port:   port,
		db:     db,
	}
}

// Start запускает HTTP сервер
func (s *Server) Start() error {
	api.Init(s.db)

	webDir := "./web"
	fileServer := http.FileServer(http.Dir(webDir))
	http.Handle("/", fileServer)

	// Лог
	s.logger.Println("The server is running on http://localhost:" + s.port)
	return http.ListenAndServe(":"+s.port, nil)
}
