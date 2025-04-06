package utils

import (
	"fmt"
	"strings"
)

func DSN(args []string) (string, error) {
	for _, arg := range args {
		if strings.HasPrefix(arg, "--dsn=") {
			return strings.TrimPrefix(arg, "--dsn="), nil
		}
	}
	return "", fmt.Errorf("--dsn argument is required")
}

func DBType(dsn string) string {
	dsn = strings.ToLower(dsn)

	switch {
	case strings.HasPrefix(dsn, "postgres://"), strings.HasPrefix(dsn, "postgresql://"):
		return "postgres"
	case strings.HasSuffix(dsn, ".db"), strings.HasSuffix(dsn, ".sqlite"), strings.HasSuffix(dsn, ".sqlite3"):
		return "sqlite"
	default:
		return "unknown"
	}
}
