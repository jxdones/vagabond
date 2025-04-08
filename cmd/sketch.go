package cmd

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/jxdones/vagabond/cmd/utils"
	"github.com/jxdones/vagabond/internal/db"
	"github.com/jxdones/vagabond/internal/schema"
)

func SketchSchema(args []string) error {
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

	var path string
	if len(args) > 0 && !strings.Contains(args[0], dsn) {
		path = args[0]
	}

	if path == "" {
		path = migrationPath
	}

	schemaPath := filepath.Join(path, "schema.sql")

	driver, err := db.New(db.Config{Type: dbType, DSN: dsn})
	if err != nil {
		return err
	}
	defer driver.Close()

	if err := schema.DumpSchema(driver, schemaPath); err != nil {
		log.Fatalf("Migration error: %v", err)
	}

	return nil
}