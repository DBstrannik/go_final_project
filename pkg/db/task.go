package db

import (
	"errors"
	"fmt"
)

// Task структура описывает одну задачу с основными полями
type Task struct {
	ID      string `json:"id,omitempty"` // ID задачи, строкой, может быть пустым при создании
	Date    string `json:"date"`         // дата задачи в формате "YYYYMMDD"
	Title   string `json:"title"`        // заголовок задачи
	Comment string `json:"comment"`      // дополнительное описание/комментарий
	Repeat  string `json:"repeat"`       // правило повтора задачи (например, "d 7" — повторять каждые 7 дней)
}

// AddTask добавляет новую задачу в таблицу scheduler и возвращает ID новой записи
func AddTask(task *Task) (int64, error) {
	if DB == nil {
		return 0, errors.New("database not initialized") // проверяем, что база инициализирована
	}

	query := `
	INSERT INTO scheduler (date, title, comment, repeat)
	VALUES (?, ?, ?, ?)
	`
	// выполняем SQL запрос на вставку новой задачи
	res, err := DB.Exec(query, task.Date, task.Title, task.Comment, task.Repeat)
	if err != nil {
		return 0, err
	}

	return res.LastInsertId() // возвращаем ID вставленной записи
}

// Tasks возвращает до limit задач, отсортированных по дате по возрастанию
func Tasks(limit int) ([]*Task, error) {
	if DB == nil {
		return nil, errors.New("database not initialized")
	}

	query := `
		SELECT id, date, title, comment, repeat
		FROM scheduler
		ORDER BY date ASC
		LIMIT ?
	`

	rows, err := DB.Query(query, limit)
	if err != nil {
		return nil, fmt.Errorf("query error: %w", err)
	}
	defer rows.Close()

	var tasks []*Task

	for rows.Next() {
		var t Task
		var id int
		var date int

		// считываем данные из строки результата в переменные
		if err := rows.Scan(&id, &date, &t.Title, &t.Comment, &t.Repeat); err != nil {
			return nil, fmt.Errorf("scan error: %w", err)
		}

		// преобразуем числовой id и дату в строковый формат
		t.ID = fmt.Sprintf("%d", id)
		t.Date = fmt.Sprintf("%08d", date)

		tasks = append(tasks, &t)
	}

	// проверяем наличие ошибок во время обхода строк
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("row iteration error: %w", err)
	}

	// если задач нет — возвращаем пустой срез, а не nil
	if tasks == nil {
		tasks = []*Task{}
	}

	return tasks, nil
}

// GetTask возвращает задачу по её ID
func GetTask(id string) (*Task, error) {
	if DB == nil {
		return nil, errors.New("database not initialized")
	}

	query := `SELECT id, date, title, comment, repeat FROM scheduler WHERE id = ?`
	var t Task
	var idNum int
	var date int

	// получаем одну строку с заданным ID
	err := DB.QueryRow(query, id).Scan(&idNum, &date, &t.Title, &t.Comment, &t.Repeat)
	if err != nil {
		return nil, err // если задача не найдена или ошибка - возвращаем ошибку
	}

	// конвертируем числовые значения в строки
	t.ID = fmt.Sprintf("%d", idNum)
	t.Date = fmt.Sprintf("%08d", date)
	return &t, nil
}

// UpdateTask обновляет поля задачи по ID
func UpdateTask(task *Task) error {
	if DB == nil {
		return errors.New("database not initialized")
	}

	query := `
	UPDATE scheduler 
	SET date = ?, title = ?, comment = ?, repeat = ?
	WHERE id = ?
	`
	// выполняем запрос обновления
	res, err := DB.Exec(query, task.Date, task.Title, task.Comment, task.Repeat, task.ID)
	if err != nil {
		return err
	}

	count, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if count == 0 {
		return fmt.Errorf("Задача не найдена") // если ни одной строки не обновлено — значит задачи нет
	}

	return nil
}

// DeleteTask удаляет задачу по ID
func DeleteTask(id string) error {
	if DB == nil {
		return errors.New("database not initialized")
	}

	query := `DELETE FROM scheduler WHERE id = ?`
	res, err := DB.Exec(query, id)
	if err != nil {
		return err
	}

	count, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if count == 0 {
		return fmt.Errorf("Задача не найдена") // если ни одной строки не удалено — значит задачи нет
	}

	return nil
}

// UpdateDate обновляет только дату задачи по ID
func UpdateDate(next string, id string) error {
	if DB == nil {
		return errors.New("database not initialized")
	}

	query := `UPDATE scheduler SET date = ? WHERE id = ?`
	res, err := DB.Exec(query, next, id)
	if err != nil {
		return err
	}

	count, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if count == 0 {
		return fmt.Errorf("Задача не найдена")
	}

	return nil
}
