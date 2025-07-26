package db

import (
	"database/sql"
	"log"
	"os"

	_ "modernc.org/sqlite"
)

// DB представляет соединение с базой данных,
// Исправлена бд, с глобальной на срукт, как в задании с посылками...
type DB struct {
	conn *sql.DB
}

const schema = `
CREATE TABLE scheduler (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	date CHAR(8) NOT NULL DEFAULT '',
	title VARCHAR(128) NOT NULL DEFAULT '',
	comment TEXT NOT NULL DEFAULT '',
	repeat VARCHAR(128) NOT NULL DEFAULT ''
);
CREATE INDEX idx_scheduler_date ON scheduler(date);
`

// Init инициализирует базу данных
func Init(dbFile string, logger *log.Logger) (*DB, error) {
	_, err := os.Stat(dbFile)
	install := false
	if err != nil {
		install = true
	}

	conn, err := sql.Open("sqlite", dbFile)
	if err != nil {
		return nil, err
	}

	if install {
		logger.Println("Create a database and a scheduler table...")
		_, err = conn.Exec(schema)
		if err != nil {
			return nil, err
		}
	}

	return &DB{conn: conn}, nil
}

// Close закрывает соединение с базой данных
func (d *DB) Close() error {
	return d.conn.Close()
}

// Conn возвращает соединение с базой данных
func (d *DB) Conn() *sql.DB {
	return d.conn
}
