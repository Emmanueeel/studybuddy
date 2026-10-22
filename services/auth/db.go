package main

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

func InitDB(dbUrl string) *sql.DB {
	db, err := sql.Open("postgres", dbUrl)
	if err != nil {
		log.Fatal("open db:", err)
	}
	if err := db.Ping(); err != nil {
		log.Fatal("ping db:", err)
	}

	schema := `
    CREATE TABLE IF NOT EXISTS users (
        id         SERIAL PRIMARY KEY,
        email      VARCHAR(255) UNIQUE NOT NULL,
        password   VARCHAR(255) NOT NULL,
        role       VARCHAR(50) NOT NULL DEFAULT 'student',
        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
    );`
	if _, err := db.Exec(schema); err != nil {
		log.Fatal("migrate:", err)
	}
	return db
}
