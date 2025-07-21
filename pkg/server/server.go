package server

import (
	"go_final_project/pkg/api"
	"log"
	"net/http"
	"os"
)

// Server структура содержит настройки сервера: логгер и порт
type Server struct {
	logger *log.Logger
	port   string
}

// NewServer создает новый сервер с логгером и портом из переменной окружения TODO_PORT,
// если переменная не задана, используется порт по умолчанию 7540
func NewServer(logger *log.Logger) *Server {
	port := os.Getenv("TODO_PORT")
	if port == "" {
		port = "7540"
	}

	return &Server{
		logger: logger,
		port:   port,
	}
}

// Start запускает HTTP сервер
func (s *Server) Start() error {
	// инициализируем API обработчики (роуты)
	api.Init()

	// директория, откуда будут обслуживаться статические файлы фронтенда
	webDir := "./web"

	// создаем HTTP файловый сервер, который будет отдавать файлы из webDir
	fileServer := http.FileServer(http.Dir(webDir))

	// подключаем файловый сервер к корневому пути "/"
	http.Handle("/", fileServer)

	// логируем информацию о запуске сервера
	s.logger.Println("The server is running on http://localhost:" + s.port)

	// запускаем HTTP сервер на заданном порту, nil — значит использовать стандартный DefaultServeMux
	return http.ListenAndServe(":"+s.port, nil)
}
