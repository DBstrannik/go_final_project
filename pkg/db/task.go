package db

import (
	"database/sql"
	"errors"
	"fmt"
)

const dateFormat = "20060102"

// Task представляет структуру задачи
type Task struct {
	ID      string `json:"id,omitempty"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}

// AddTask добавляет новую задачу в базу данных
func (d *DB) AddTask(task *Task) (int64, error) {
	query := `
	INSERT INTO scheduler (date, title, comment, repeat)
	VALUES (?, ?, ?, ?)
	`
	res, err := d.conn.Exec(query, task.Date, task.Title, task.Comment, task.Repeat)
	if err != nil {
		return 0, fmt.Errorf("database error: %w", err)
	}
	return res.LastInsertId()
}

// Tasks возвращает список задач с ограничением по количеству
func (d *DB) Tasks(limit int) ([]*Task, error) {
	query := `
		SELECT id, date, title, comment, repeat
		FROM scheduler
		ORDER BY date ASC
		LIMIT ?
	`

	rows, err := d.conn.Query(query, limit)
	if err != nil {
		return nil, fmt.Errorf("query error: %w", err)
	}
	defer rows.Close()

	var tasks []*Task

	for rows.Next() {
		var t Task
		var id int
		var date int

		if err := rows.Scan(&id, &date, &t.Title, &t.Comment, &t.Repeat); err != nil {
			return nil, fmt.Errorf("scan error: %w", err)
		}

		t.ID = fmt.Sprintf("%d", id)
		t.Date = fmt.Sprintf("%08d", date)

		tasks = append(tasks, &t)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("row iteration error: %w", err)
	}

	if tasks == nil {
		tasks = []*Task{}
	}

	return tasks, nil
}

// GetTask возвращает задачу по ID
func (d *DB) GetTask(id string) (*Task, error) {
	query := `SELECT id, date, title, comment, repeat FROM scheduler WHERE id = ?`
	var t Task
	var idNum int
	var date int

	err := d.conn.QueryRow(query, id).Scan(&idNum, &date, &t.Title, &t.Comment, &t.Repeat)
	if err != nil {
		if err == sql.ErrNoRows {
			// Исправлено: Ошибка на английском и с маленькой буквы
			return nil, errors.New("task not found")
		}
		return nil, err
	}

	t.ID = fmt.Sprintf("%d", idNum)
	t.Date = fmt.Sprintf("%08d", date)
	return &t, nil
}

// UpdateTask обновляет существующую задачу
func (d *DB) UpdateTask(task *Task) error {
	query := `
	UPDATE scheduler 
	SET date = ?, title = ?, comment = ?, repeat = ?
	WHERE id = ?
	`
	res, err := d.conn.Exec(query, task.Date, task.Title, task.Comment, task.Repeat, task.ID)
	if err != nil {
		return fmt.Errorf("database error: %w", err)
	}

	count, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("rows affected error: %w", err)
	}
	if count == 0 {
		// Исправлено: Ошибка на английском и с маленькой буквы
		return errors.New("task not found")
	}

	return nil
}

// DeleteTask удаляет задачу по ID
func (d *DB) DeleteTask(id string) error {
	query := `DELETE FROM scheduler WHERE id = ?`
	res, err := d.conn.Exec(query, id)
	if err != nil {
		return fmt.Errorf("database error: %w", err)
	}

	count, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("rows affected error: %w", err)
	}
	if count == 0 {
		// Исправлено: Ошибка на английском и с маленькой буквы
		return errors.New("task not found")
	}

	return nil
}

// UpdateDate обновляет дату выполнения задачи
func (d *DB) UpdateDate(next string, id string) error {
	query := `UPDATE scheduler SET date = ? WHERE id = ?`
	res, err := d.conn.Exec(query, next, id)
	if err != nil {
		return fmt.Errorf("database error: %w", err)
	}

	count, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("rows affected error: %w", err)
	}
	if count == 0 {
		// Исправлено: Ошибка на английском и с маленькой буквы
		return errors.New("task not found")
	}

	return nil
}
