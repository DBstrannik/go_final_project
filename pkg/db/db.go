package db

import (
	"database/sql"
	"log"
	"os"

	_ "modernc.org/sqlite" // SQLite драйвер для Go
)

var DB *sql.DB // глобальная переменная с подключением к базе данных

// SQL-схема для создания таблицы scheduler и индекса по дате
const schema = `
CREATE TABLE scheduler (
	id INTEGER PRIMARY KEY AUTOINCREMENT,      -- уникальный ID задачи, автоинкремент
	date CHAR(8) NOT NULL DEFAULT '',           -- дата задачи в формате "YYYYMMDD"
	title VARCHAR(128) NOT NULL DEFAULT '',     -- заголовок задачи
	comment TEXT NOT NULL DEFAULT '',            -- комментарий/описание задачи
	repeat VARCHAR(128) NOT NULL DEFAULT ''     -- правило повтора задачи (например "d 7")
);
CREATE INDEX idx_scheduler_date ON scheduler(date);  -- индекс по дате для ускорения выборок
`

// Init инициализирует базу данных:
// - проверяет, существует ли файл базы данных
// - если нет, создаёт новую базу и нужные таблицы
// - открывает соединение и сохраняет его в глобальной переменной DB
func Init(dbFile string, logger *log.Logger) error {
	// Проверяем наличие файла базы данных
	_, err := os.Stat(dbFile)
	install := false
	if err != nil {
		install = true // файл отсутствует — нужно создать базу с таблицей
	}

	// Открываем соединение с SQLite базой по пути dbFile
	DB, err = sql.Open("sqlite", dbFile)
	if err != nil {
		return err
	}

	// Если база новая, создаём таблицу и индекс
	if install {
		logger.Println("Create a database and a scheduler table...")
		_, err = DB.Exec(schema)
		if err != nil {
			return err
		}
	}

	return nil
}
