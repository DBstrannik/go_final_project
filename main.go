package main

import (
	"go_final_project/pkg/db"
	"go_final_project/pkg/server"
	"log"
	"os"
)

func main() {
	// Создаём логгер, который будет писать в консоль с префиксом "INFO:" и меткой времени
	logger := log.New(os.Stdout, "INFO: ", log.LstdFlags)

	// Подключаемся к SQLite-базе данных (файл scheduler.db).
	// Если база не существует, скорее всего будет создана.
	if err := db.Init("scheduler.db", logger); err != nil {
		// Прерываем выполнение, если не удалось инициализировать базу
		logger.Fatal("error initializing database:", err)
	}

	// Создаём HTTP-сервер и передаём в него логгер для последующего логирования внутри сервера
	srv := server.NewServer(logger)

	// Запускаем сервер — например, он может слушать порт и обрабатывать запросы
	if err := srv.Start(); err != nil {
		// Если порт занят или сервер не может стартануть, приложение аварийно завершится
		logger.Fatal("Error starting server: ", err)
	}
}
