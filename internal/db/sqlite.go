package db

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

type SQLite struct {
	conn *sql.DB
}

func (s *SQLite) Connect(dsn string) error {
	var err error
	conn, err := sql.Open("sqlite3", dsn)
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	s.conn = conn
	return s.createMigrationsTable()
}

func (s *SQLite) Close() error {
	if s.conn == nil {
		return fmt.Errorf("no open connection to close")
	}
	return s.conn.Close()
}

func (s *SQLite) createMigrationsTable() error {
	query := `
	CREATE TABLE IF NOT EXISTS vagabond_migrations (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		migration_id TEXT NOT NULL UNIQUE,
		applied_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	)`
	_, err := s.conn.Exec(query)
	return err
}

func (s *SQLite) GetAppliedMigrations() (map[string]bool, error) {
	rows, err := s.conn.Query("SELECT migration_id FROM vagabond_migrations ORDER BY applied_at ASC")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	applied := make(map[string]bool)
	for rows.Next() {
		var migration_id string
		if err := rows.Scan(&migration_id); err != nil {
			return nil, err
		}
		applied[migration_id] = true
	}
	return applied, nil
}

func (s *SQLite) GetAppliedMigrationsList() ([]string, error) {
	rows, err := s.conn.Query("SELECT migration_id FROM vagabond_migrations ORDER BY applied_at ASC")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []string
	for rows.Next() {
		var migration_id string
		if err := rows.Scan(&migration_id); err != nil {
			return nil, err
		}
		list = append(list, migration_id)
	}
	return list, nil
}

func (s *SQLite) ExecuteMigration(filePath string) error {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("error reading file %s: %w", filePath, err)
	}

	tx, err := s.conn.Begin()
	if err != nil {
		return err
	}

	if _, err := tx.Exec(string(data)); err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to execute migration: %w", err)
	}

	_, migration_id := filepath.Split(filePath)
	_, err = tx.Exec("INSERT INTO vagabond_migrations (migration_id) VALUES (?)", migration_id)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to record migration: %w", err)
	}

	return tx.Commit()
}

func (s *SQLite) RollbackMigration(filePath string, name string) error {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("error reading file %s: %w", filePath, err)
	}

	tx, err := s.conn.Begin()
	if err != nil {
		return err
	}

	if _, err := tx.Exec(string(data)); err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to execute migration: %w", err)
	}

	_, err = tx.Exec("DELETE from vagabond_migrations WHERE name = ?", name)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to delete migration record: %w", err)
	}

	return tx.Commit()
}

func (s *SQLite) DumpSchema() (string, error) {
	var schema strings.Builder

	schema.WriteString("-- This file has been automatically generated based on the current database state.\n")
	schema.WriteString("-- Manual modification of this file is not recommended. Use database migrations for schema changes.\n\n")
	schema.WriteString("PRAGMA foreign_keys = OFF;\n\n")

	rows, err := s.conn.Query(`
		SELECT sql FROM sqlite_master
		WHERE type = 'table' AND name NOT LIKE 'sqlite_%' AND sql IS NOT NULL
	`)
	if err != nil {
		return "", fmt.Errorf("failed to fetch schema: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var stmt string
		if err := rows.Scan(&stmt); err != nil {
			return "", err
		}
		schema.WriteString(stmt + ";\n\n")
	}

	schema.WriteString("PRAGMA foreign_keys = ON;\n\n")

	migrationRows, err := s.conn.Query(`SELECT id, migration_id FROM vagabond_migrations ORDER BY id`)
	if err != nil {
		return "", fmt.Errorf("failed to fetch applied migrations: %w", err)
	}
	defer migrationRows.Close()

	var values []string
	for migrationRows.Next() {
		var id int
		var migration_id string
		if err := migrationRows.Scan(&id, &migration_id); err != nil {
			return "", err
		}
		values = append(values, fmt.Sprintf("(%d, '%s')", id, migration_id))
	}

	if len(values) > 0 {
		schema.WriteString("INSERT INTO vagabond_migrations (id, migration_id) VALUES\n\t")
		schema.WriteString(strings.Join(values, ",\n\t"))
		schema.WriteString(";\n")
	}

	return schema.String(), nil
}