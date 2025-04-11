package commands

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/jxdones/vagabond/commands/utils"
	"github.com/jxdones/vagabond/internal/db"
	"github.com/jxdones/vagabond/internal/migrations"
)

const defaultRollbackCount = 1

func UnpackMigrations(args []string) error {
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

	var n int
	if len(args) > 0 && !strings.Contains(args[0], dsn) {
		parsed, err := strconv.Atoi(args[0])
		if err != nil {
			return fmt.Errorf("invalid option: provide a number")
		}
		n = parsed
	}

	if n == 0 {
		n = defaultRollbackCount
	}

	driver, err := db.New(db.Config{Type: dbType, DSN: dsn})
	if err != nil {
		return err
	}
	defer driver.Close()

	return migrations.RollbackMigrations(driver, n)
}
