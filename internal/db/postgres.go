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
	var schema strings.Builder

	schema.WriteString("-- This file has been automatically generated based on the current database state.\n")
	schema.WriteString("-- Manual modification of this file is not recommended. Use database migrations for schema changes.\n\n")

	tables, err := p.getTables()
	if err != nil {
		return "", err
	}

	for _, table := range tables {
		createStmt, err := p.getCreateTableStmt(table)
		if err != nil {
			return "", err
		}
		schema.WriteString(createStmt + ";\n\n")
	}

	indexes, err := p.getIndexes()
	if err != nil {
		return "", err
	}
	for _, idx := range indexes {
		schema.WriteString(idx + ";\n\n")
	}

	enums, err := p.getEnums()
	if err != nil {
		return "", err
	}
	for _, enum := range enums {
		schema.WriteString(enum + ";\n\n")
	}

	migrationRows, err := p.conn.Query(`SELECT id, migration_id FROM vagabond_migrations ORDER BY id`)
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

	return strings.TrimSpace(schema.String()), nil
}

func (p *Postgres) getTables() ([]string, error) {
	rows, err := p.conn.Query(`
		SELECT tablename
		FROM pg_tables
		WHERE schemaname = 'public'
		ORDER BY tablename
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tables []string
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return nil, err
		}
		if name == "_vagabond_migrations" {
			continue
		}
		tables = append(tables, name)
	}
	return tables, nil
}

func (p *Postgres) getCreateTableStmt(table string) (string, error) {
	columns, err := p.getColumns(table)
	if err != nil {
		return "", err
	}

	constraints, err := p.getConstraints(table)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("CREATE TABLE %q (\n%s\n%s)", table, strings.Join(columns, ",\n"), constraints), nil
}

func (p *Postgres) getColumns(table string) ([]string, error) {
	rows, err := p.conn.Query(`
		SELECT column_name, data_type, is_nullable, column_default
		FROM information_schema.columns
		WHERE table_schema = 'public' AND table_name = $1
		ORDER BY ordinal_position
	`, table)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var cols []string
	for rows.Next() {
		var name, dataType, nullable string
		var defaultVal *string

		if err := rows.Scan(&name, &dataType, &nullable, &defaultVal); err != nil {
			return nil, err
		}

		col := fmt.Sprintf("  %q %s", name, dataType)
		if defaultVal != nil {
			col += fmt.Sprintf(" DEFAULT %s", *defaultVal)
		}
		if nullable == "NO" {
			col += " NOT NULL"
		}
		cols = append(cols, col)
	}
	return cols, nil
}

func (p *Postgres) getConstraints(table string) (string, error) {
	rows, err := p.conn.Query(`
		SELECT
			c.constraint_type,
			c.constraint_name,
			string_agg(cols.column_name, ',') AS columns,
			fk.foreign_table_name,
			string_agg(fk.foreign_column_name, ',') AS foreign_columns
		FROM information_schema.table_constraints c
		JOIN information_schema.constraint_column_usage cols
			ON c.constraint_name = cols.constraint_name
		LEFT JOIN (
			SELECT
				tc.constraint_name,
				ccu.table_name AS foreign_table_name,
				ccu.column_name AS foreign_column_name
			FROM information_schema.referential_constraints rc
			JOIN information_schema.table_constraints tc
				ON rc.constraint_name = tc.constraint_name
			JOIN information_schema.constraint_column_usage ccu
				ON rc.unique_constraint_name = ccu.constraint_name
		) fk ON c.constraint_name = fk.constraint_name
		WHERE c.table_name = $1 AND c.table_schema = 'public'
		GROUP BY c.constraint_type, c.constraint_name, fk.foreign_table_name
	`, table)
	if err != nil {
		return "", err
	}
	defer rows.Close()

	var constraints []string
	for rows.Next() {
		var ctype, name, cols, ftable, fcols string
		_ = rows.Scan(&ctype, &name, &cols, &ftable, &fcols)

		switch ctype {
		case "PRIMARY KEY":
			constraints = append(constraints, fmt.Sprintf("  ,PRIMARY KEY (%s)", quoteIdentList(cols)))
		case "UNIQUE":
			constraints = append(constraints, fmt.Sprintf("  ,UNIQUE (%s)", quoteIdentList(cols)))
		case "FOREIGN KEY":
			constraints = append(constraints, fmt.Sprintf("  ,FOREIGN KEY (%s) REFERENCES %s (%s)",
				quoteIdentList(cols), ftable, quoteIdentList(fcols)))
		}
	}
	return strings.Join(constraints, "\n"), nil
}

func quoteIdentList(csv string) string {
	parts := strings.Split(csv, ",")
	for i, p := range parts {
		parts[i] = fmt.Sprintf(`"%s"`, strings.TrimSpace(p))
	}
	return strings.Join(parts, ", ")
}

func (p *Postgres) getEnums() ([]string, error) {
	rows, err := p.conn.Query(`
		SELECT t.typname, string_agg(e.enumlabel, ',' ORDER BY e.enumsortorder)
		FROM pg_type t
		JOIN pg_enum e ON t.oid = e.enumtypid
		JOIN pg_namespace n ON n.oid = t.typnamespace
		WHERE n.nspname = 'public'
		GROUP BY t.typname
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var enums []string
	for rows.Next() {
		var name, values string
		_ = rows.Scan(&name, &values)

		vals := strings.Split(values, ",")
		for i := range vals {
			vals[i] = fmt.Sprintf("'%s'", vals[i])
		}
		enums = append(enums, fmt.Sprintf("CREATE TYPE %q AS ENUM (%s)", name, strings.Join(vals, ", ")))
	}
	return enums, nil
}

func (p *Postgres) getIndexes() ([]string, error) {
	rows, err := p.conn.Query(`
		SELECT indexname, indexdef
		FROM pg_indexes
		WHERE schemaname = 'public'
		AND indexdef NOT ILIKE '%pkey%'
		AND indexdef NOT ILIKE '%unique%'
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var indexes []string
	for rows.Next() {
		var name, def string
		if err := rows.Scan(&name, &def); err != nil {
			return nil, err
		}
		indexes = append(indexes, def)
	}
	return indexes, nil
}
