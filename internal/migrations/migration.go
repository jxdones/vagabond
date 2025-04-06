package migrations

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/jxdones/vagabond/internal/db"
)

const migrationsPath = "migrations"

func CreateMigration(name string) error {
	timestamp := time.Now().Format("20060102150405")
	upFileName := fmt.Sprintf("%s_%s_up.sql", timestamp, name)
	downFileName := fmt.Sprintf("%s_%s_down.sql", timestamp, name)

	upFilePath := fmt.Sprintf("%s/%s", migrationsPath, upFileName)
	downFilePath := fmt.Sprintf("%s/%s", migrationsPath, downFileName)

	upFile, err := os.Create(upFilePath)
	if err != nil {
		return fmt.Errorf("error creating up migration file: %w", err)
	}
	defer upFile.Close()

	downFile, err := os.Create(downFilePath)
	if err != nil {
		return fmt.Errorf("error creating down migration file: %w", err)
	}
	defer downFile.Close()

	upFile.WriteString(fmt.Sprintf("-- %s\n--  Write your SQL to apply this migration.\n", upFileName))
	downFile.WriteString(fmt.Sprintf("-- %s\n-- Write your SQL to rollback this migration.\n", downFileName))
	fmt.Printf("%s created\n", upFilePath)
	fmt.Printf("%s created\n", downFilePath)
	return nil
}

func ApplyMigrations(driver db.Driver, dbType string) error {
	applied, err := driver.GetAppliedMigrations()
	if err != nil {
		return fmt.Errorf("could not get applied migrations: %w", err)
	}

	files, err := filepath.Glob(filepath.Join(migrationsPath, "*_up.sql"))
	if err != nil {
		return fmt.Errorf("failed to list migration files: %w", err)
	}
	if len(files) == 0 {
		fmt.Println("No migrations found.")
		return nil
	}

	sort.Strings(files)
	var pending []string
	for _, f := range files {
		_, name := filepath.Split(f)
		if !applied[name] {
			pending = append(pending, f)
		}
	}

	if len(pending) == 0 {
		fmt.Println("No new migrations to apply.")
		return nil
	}

	for _, file := range pending {
		fmt.Printf("Applying migration: %s\n", filepath.Base(file))
		if err := driver.ExecuteMigration(file); err != nil {
			return fmt.Errorf("error applying %s: %w", file, err)
		}
	}

	fmt.Printf("Applied %d migration(s).\n", len(pending))
	return nil
}

func RollbackMigrations(driver db.Driver, n int) error {
	appliedMigrations, err := driver.GetAppliedMigrationsList()
	if err != nil {
		return fmt.Errorf("could not get applied migrations: %w", err)
	}

	total := len(appliedMigrations)
	if total == 0 {
		return nil
	}

	// ensure that n will always be capped to total 
	if n > total {
		n = total
	}

	toRollback := appliedMigrations[total-n:]
	for i := len(toRollback) - 1; i >= 0; i-- {
		name := toRollback[i]
		downFile := downFileName(name)
		path := filepath.Join(migrationsPath, downFile)

		fmt.Printf("Rolling back migration: %s\n", name)
		if err := driver.RollbackMigration(path, name); err != nil {
			return fmt.Errorf("failed to rollback %s: %w", name, err)
		}
		fmt.Printf("Rolled back: %s\n", name)
	}
	fmt.Printf("Rolled back %d migration(s).\n", n)
	return nil
}

func downFileName(upFile string) string {
	if !strings.HasSuffix(upFile, "_up.sql") {
		return ""
	}
	return upFile[:len(upFile)-len("_up.sql")] + "_down.sql"
}