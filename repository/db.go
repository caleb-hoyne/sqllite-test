package db

import (
	"database/sql"
	"errors"
)

type Repository struct {
	DB *sql.DB
}

func (h *Repository) GetNameByID(id int) (string, error) {
	var name string
	err := h.DB.QueryRow("SELECT name FROM test WHERE id = ?", id).Scan(&name)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", nil
		}
		return "", err
	}
	return name, nil
}
