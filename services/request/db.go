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
	CREATE TABLE IF NOT EXISTS requests (
		id           SERIAL PRIMARY KEY,
		initiator_id INT NOT NULL,
		target_id    INT NOT NULL,
		discipline   VARCHAR(255),
		message      TEXT,
		status       VARCHAR(50) NOT NULL DEFAULT 'created',
		created_at   TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at   TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);`
	if _, err := db.Exec(schema); err != nil {
		log.Fatal("migrate:", err)
	}
	return db
}
