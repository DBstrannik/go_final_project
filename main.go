package main

import (
	"go_final_project/pkg/db"
	"go_final_project/pkg/server"
	"log"
	"os"
)

func main() {
	// Создаем логгер для вывода информации о работе программы
	logger := log.New(os.Stdout, "INFO: ", log.LstdFlags)

	// Инициализируем подключение к базе данных
	dbInstance, err := db.Init("scheduler.db", logger)
	if err != nil {
		logger.Fatal("Error connecting to the database:", err)
	}

	// Исправлено: Добавлено закрытие подключения к БД при завершении программы
	defer func() {
		if err := dbInstance.Close(); err != nil {
			logger.Printf("Error closing the DATABASE connection: %v", err)
		}
		logger.Println("The database connection is closed")
	}()

	// Создаем новый сервер
	srv := server.NewServer(logger, dbInstance)

	// Запускаем сервер
	logger.Println("Server starting...")
	if err := srv.Start(); err != nil {
		logger.Fatal("Server error:", err)
	}
}
