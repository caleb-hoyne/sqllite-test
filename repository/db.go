package db

import (
	"database/sql"
	"errors"
	"github.com/davecgh/go-spew/spew"
	"github.com/mattn/go-sqlite3"
)

var (
	ErrIDAlreadyExists = errors.New("id already exists")
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
	spew.Dump(err)
	var sqErr sqlite3.Error
	if errors.As(err, &sqErr) && sqErr.Code == sqlite3.ErrConstraint {
		return ErrIDAlreadyExists
	}
	return err
}
