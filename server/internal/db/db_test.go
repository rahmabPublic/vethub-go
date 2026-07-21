package db

import (
	"database/sql"
	"strings"
	"testing"
)

func TestOpen_InitializesSchemaAndSeed(t *testing.T) {
	t.Parallel()

	database, err := Open()
	if err != nil {
		t.Fatalf("expected Open to succeed, got error: %v", err)
	}
	defer database.Close()

	var ownerCount int
	if err := database.QueryRow("SELECT COUNT(*) FROM owners").Scan(&ownerCount); err != nil {
		t.Fatalf("failed to query owners count: %v", err)
	}
	if ownerCount == 0 {
		t.Fatalf("expected seeded owners, got 0")
	}
}

func TestRunSQL_ReturnsReadErrorForMissingFile(t *testing.T) {
	t.Parallel()

	database, err := sql.Open("sqlite", "file:test-read-error?mode=memory&cache=shared")
	if err != nil {
		t.Fatalf("failed to open sqlite db: %v", err)
	}
	defer database.Close()

	err = runSQL(database, "missing-file.sql")
	if err == nil {
		t.Fatalf("expected error for missing SQL file")
	}
	if !strings.Contains(err.Error(), "read missing-file.sql") {
		t.Fatalf("expected read error, got: %v", err)
	}
}

func TestRunSQL_ReturnsExecErrorWhenSchemaMissing(t *testing.T) {
	t.Parallel()

	database, err := sql.Open("sqlite", "file:test-exec-error?mode=memory&cache=shared")
	if err != nil {
		t.Fatalf("failed to open sqlite db: %v", err)
	}
	defer database.Close()

	err = runSQL(database, "seed.sql")
	if err == nil {
		t.Fatalf("expected exec error when running seed without schema")
	}
	if !strings.Contains(err.Error(), "exec seed.sql") {
		t.Fatalf("expected exec error, got: %v", err)
	}
}
