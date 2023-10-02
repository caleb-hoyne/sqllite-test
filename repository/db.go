package db

import (
	"database/sql"
	"errors"
	"github.com/mattn/go-sqlite3"
)

var (
	ErrIDAlreadyExists = errors.New("id already exists")
	ErrNotFound        = errors.New("id not found")
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

func (h *Repository) StoreUser(id int, name string) error {
	_, err := h.DB.Exec("INSERT INTO test VALUES (?, ?)", id, name)
	var sqErr sqlite3.Error
	if errors.As(err, &sqErr) && sqErr.Code == sqlite3.ErrConstraint {
		return ErrIDAlreadyExists
	}
	return err
}

func (h *Repository) UpdateUser(id int, newName string) error {
	_, err := h.DB.Exec("UPDATE test SET name = ? WHERE id = ?", newName, id)
	var sqErr sqlite3.Error
	if errors.As(err, &sqErr) && sqErr.Code == sqlite3.ErrNotFound {
		return h.StoreUser(id, newName)
	}
	return err
}
