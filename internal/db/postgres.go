package db

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	_ "github.com/lib/pq"
)

type Postgres struct {
	conn *sql.DB
}

func (p *Postgres) Connect(dsn string) error {
	var err error
	conn, err := sql.Open("postgres", dsn)
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	p.conn = conn
	return p.createMigrationsTable()
}

func (p *Postgres) createMigrationsTable() error {
	query := `
	CREATE TABLE IF NOT EXISTS vagabond_migrations (
		id SERIAL PRIMARY KEY,
		migration_id VARCHAR(255) NOT NULL UNIQUE,
		applied_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	)`
	_, err := p.conn.Exec(query)
	return err
}

func (p *Postgres) Close() error {
	if p.conn == nil {
		return fmt.Errorf("no open connection to close")
	}
	return p.conn.Close()
}

func (p *Postgres) GetAppliedMigrations() (map[string]bool, error) {
	rows, err := p.conn.Query("SELECT migration_id FROM vagabond_migrations ORDER BY applied_at ASC")
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

func (p *Postgres) GetAppliedMigrationsList() ([]string, error) {
	rows, err := p.conn.Query("SELECT migration_id FROM vagabond_migrations ORDER BY applied_at ASC")
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

func (p *Postgres) ExecuteMigration(filePath string) error {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("error reading file %s: %w", filePath, err)
	}

	tx, err := p.conn.Begin()
	if err != nil {
		return err
	}

	if _, err := tx.Exec(string(data)); err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to execute migration: %w", err)
	}

	_, migrationFile := filepath.Split(filePath)
	migrationID := strings.TrimSuffix(migrationFile, ".sql")
	_, err = tx.Exec("INSERT INTO vagabond_migrations (migration_id) VALUES ($1)", migrationID)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to record migration: %w", err)
	}

	return tx.Commit()
}

func (p *Postgres) RollbackMigration(filePath string, name string) error {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("error reading file %s: %w", filePath, err)
	}

	tx, err := p.conn.Begin()
	if err != nil {
		return err
	}

	if _, err := tx.Exec(string(data)); err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to execute migration: %w", err)
	}

	migrationID := strings.TrimSuffix(name, ".sql")
	_, err = tx.Exec("DELETE from vagabond_migrations WHERE migration_id = $1", migrationID)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to delete migration record: %w", err)
	}

	return tx.Commit()
}

func (p *Postgres) DumpSchema() (string, error) {
	// to be implemented
	return "", nil
}
