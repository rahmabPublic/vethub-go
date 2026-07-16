package db

import (
	"database/sql"
	"embed"
	"fmt"
	"log"

	_ "modernc.org/sqlite"
)

//go:embed schema.sql seed.sql
var sqlFiles embed.FS

func Open() (*sql.DB, error) {
	db, err := sql.Open("sqlite", "file:vethub?mode=memory&cache=shared")
	if err != nil {
		return nil, fmt.Errorf("open database: %w", err)
	}

	if _, err := db.Exec("PRAGMA foreign_keys = ON"); err != nil {
		db.Close()
		return nil, fmt.Errorf("enable foreign keys: %w", err)
	}

	if err := runSQL(db, "schema.sql"); err != nil {
		db.Close()
		return nil, fmt.Errorf("run schema: %w", err)
	}

	if err := runSQL(db, "seed.sql"); err != nil {
		db.Close()
		return nil, fmt.Errorf("run seed: %w", err)
	}

	var count int
	if err := db.QueryRow("SELECT COUNT(*) FROM owners").Scan(&count); err != nil {
		db.Close()
		return nil, fmt.Errorf("verify seed data: %w", err)
	}
	log.Printf("Database initialized: %d owners loaded", count)

	return db, nil
}

func runSQL(db *sql.DB, filename string) error {
	data, err := sqlFiles.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("read %s: %w", filename, err)
	}
	if _, err := db.Exec(string(data)); err != nil {
		return fmt.Errorf("exec %s: %w", filename, err)
	}
	return nil
}
