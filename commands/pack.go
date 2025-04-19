package commands

import (
	"fmt"
	"log"
	"os"

	"github.com/jxdones/vagabond/commands/utils"
	"github.com/jxdones/vagabond/internal/db"
	"github.com/jxdones/vagabond/internal/migrations"
)

func PackMigration(args []string) error {
	dsn, err := utils.DSN(args)
	if err != nil {
		return err
	}

	dbType := utils.DBType(dsn)
	if dbType == "unknown" {
		return fmt.Errorf("could not determine database type from DSN")
	}

	if _, err := os.Stat(migrationPath); os.IsNotExist(err) {
		return fmt.Errorf("missing migrations directory")
	}
	driver, err := db.New(db.Config{Type: dbType, DSN: dsn})
	if err != nil {
		return err
	}
	defer driver.Close()

	if err := migrations.ApplyMigrations(driver, dbType); err != nil {
		return fmt.Errorf("error applying migrations: %w", err)
	}

	return nil
}
