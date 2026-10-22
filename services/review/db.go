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
	CREATE TABLE IF NOT EXISTS completed_requests (
		request_id    INT PRIMARY KEY,
		initiator_id  INT NOT NULL,
		target_id     INT NOT NULL
	);
	CREATE TABLE IF NOT EXISTS reviews (
		id         SERIAL PRIMARY KEY,
		request_id INT NOT NULL,
		author_id  INT NOT NULL,
		target_id  INT NOT NULL,
		score      INT NOT NULL CHECK (score BETWEEN 1 AND 5),
		comment    TEXT,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		UNIQUE (request_id, author_id)
	);`
	if _, err := db.Exec(schema); err != nil {
		log.Fatal("migrate:", err)
	}
	return db
}
