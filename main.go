package main

import (
	"database/sql"
	"github.com/caleb-hoyne/sqllite-test/handler"
	"github.com/caleb-hoyne/sqllite-test/repository"
	_ "github.com/mattn/go-sqlite3"
	"net/http"
)

const dbPath = "test.db"

func main() {
	pool, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		panic(err)
	}
	defer pool.Close()

	_, err = pool.Exec("CREATE TABLE IF NOT EXISTS test (id INTEGER PRIMARY KEY, name TEXT)")
	if err != nil {
		panic(err)
	}

	//_, err = pool.Exec("INSERT INTO test VALUES (1, 'test')")
	//if err != nil {
	//	panic(err)
	//}

	s := http.Server{
		Addr: ":8080",
		Handler: &handler.RequestHandler{
			R: &db.Repository{
				DB: pool,
			},
		},
	}
	panic(s.ListenAndServe())
}
