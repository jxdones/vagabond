package cmd

import (
	"fmt"
	"os"

	"github.com/jxdones/vagabond/internal/migrations"
)
const migrationPath = "migrations"

func Create(args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("migration name required")
	}
	name := args[0]

	if _, err := os.Stat(migrationPath); os.IsNotExist(err) {
		err := os.Mkdir(migrationPath, 0755)
		if err != nil {
			fmt.Println("error creating migrations directory:", err)
		}
	}

	err := migrations.CreateMigration(name)
	if err != nil {
		return err
	}
	return nil
}