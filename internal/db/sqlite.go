package db

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

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
	_, err := s.conn.Exec(`
		CREATE TABLE IF NOT EXISTS vagabond_migrations (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL UNIQUE,
			applied_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)
	`)
	return err
}

func (s *SQLite) GetAppliedMigrations() (map[string]bool, error) {
	rows, err := s.conn.Query("SELECT name FROM vagabond_migrations ORDER BY applied_at ASC")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	applied := make(map[string]bool)
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return nil, err
		}
		applied[name] = true
	}
	return applied, nil
}

func (s *SQLite) GetAppliedMigrationsList() ([]string, error) {
	rows, err := s.conn.Query("SELECT name FROM vagabond_migrations ORDER BY applied_at ASC")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []string
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return nil, err
		}
		list = append(list, name)
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

	_, name := filepath.Split(filePath)
	_, err = tx.Exec("INSERT INTO vagabond_migrations (name) VALUES (?)", name)
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