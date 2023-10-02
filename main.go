package main

import (
	"database/sql"
	"github.com/caleb-hoyne/slogctx"
	"github.com/caleb-hoyne/sqllite-test/handler"
	"github.com/caleb-hoyne/sqllite-test/repository"
	_ "github.com/mattn/go-sqlite3"
	"log/slog"
	"net/http"
	"os"
)

const dbPath = "test.db"

func main() {
	pool, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		panic(err)
	}
	defer pool.Close()

	slog.SetDefault(slog.New(&slogctx.Handler{
		Handler: slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{}),
	}))
	_, err = pool.Exec("CREATE TABLE IF NOT EXISTS test (id INTEGER PRIMARY KEY, name TEXT)")
	if err != nil {
		panic(err)
	}

	nameHandler := &handler.RequestHandler{
		R: &db.Repository{
			DB: pool,
		},
	}

	mux := http.NewServeMux()
	mux.Handle("/name/", nameHandler)

	s := http.Server{
		Addr:    ":8080",
		Handler: mux,
	}
	panic(s.ListenAndServe())
}
