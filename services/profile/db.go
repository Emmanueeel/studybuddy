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
    CREATE TABLE IF NOT EXISTS profiles (
        user_id     INT PRIMARY KEY,
        name        VARCHAR(255),
        disciplines TEXT[],
        skills      TEXT[],
        level       VARCHAR(50),
        schedule    TEXT[],
        format      VARCHAR(50),
        rating      REAL DEFAULT 0,
        updated_at  TIMESTAMP DEFAULT CURRENT_TIMESTAMP
    );`
	if _, err := db.Exec(schema); err != nil {
		log.Fatal("migrate:", err)
	}
	return db
}
