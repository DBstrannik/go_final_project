package api

import (
	"errors"
	"strconv"
	"strings"
	"time"
)

const DateFormat = "20060102" // формат даты без времени, например "20250721"

// afterNow сравнивает две даты, игнорируя время, возвращая true если date > now
func afterNow(date, now time.Time) bool {
	dy, dm, dd := date.Date()
	ny, nm, nd := now.Date()

	if dy != ny {
		return dy > ny
	}
	if dm != nm {
		return dm > nm
	}
	return dd > nd
}

// NextDate рассчитывает следующую дату выполнения задачи в зависимости от правила повторения
func NextDate(now time.Time, dstart string, repeat string) (string, error) {
	// если правило повторения пустое, задача считается не повторяющейся (удаляется)
	if repeat == "" {
		return "", errors.New("no repeat rule: task will be deleted")
	}

	// парсим стартовую дату из строки
	startDate, err := time.Parse(DateFormat, dstart)
	if err != nil {
		return "", errors.New("invalid start date")
	}

	// разбиваем правило повторения, например "d 3" (каждые 3 дня)
	parts := strings.Split(repeat, " ")

	switch parts[0] {
	case "d": // повторение по дням
		if len(parts) != 2 {
			return "", errors.New("invalid d rule format")
		}
		days, err := strconv.Atoi(parts[1])
		if err != nil || days < 1 || days > 400 {
			return "", errors.New("invalid number of days")
		}
		// прибавляем дни пока дата не станет строго позже текущей (now)
		for {
			startDate = startDate.AddDate(0, 0, days)
			if afterNow(startDate, now) {
				break
			}
		}
		return startDate.Format(DateFormat), nil

	case "y": // повторение по годам (ежегодно)
		// прибавляем по одному году пока дата не станет строго позже now
		for {
			startDate = startDate.AddDate(1, 0, 0)
			if afterNow(startDate, now) {
				break
			}
		}
		return startDate.Format(DateFormat), nil

	default:
		// если правило не поддерживается — возвращаем ошибку
		return "", errors.New("unsupported repeat rule")
	}
}
