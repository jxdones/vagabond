package schema

import (
	"fmt"
	"os"

	"github.com/jxdones/vagabond/internal/db"
)

func DumpSchema(driver db.Driver, outputPath string) error {
	schema, err := driver.DumpSchema()
	if err != nil {
		return fmt.Errorf("failed to dump schema: %w", err)
	}

	if err := os.WriteFile(outputPath, []byte(schema), 0o644); err != nil {
		return fmt.Errorf("failed to create scheme file: %w", err)
	}
	return nil
}
