package db

import (
	"fmt"
	"strings"
)

type Config struct {
	Type string
	DSN  string
}

func New(cfg Config) (Driver, error) {
	var driver Driver
	switch strings.ToLower(cfg.Type) {
	case "sqlite":
		driver = &SQLite{}
	case "postgres":
		driver = &Postgres{}
	default:
		return nil, fmt.Errorf("unsupported database: %s", cfg.Type)
	}

	if err := driver.Connect(cfg.DSN); err != nil {
		return nil, fmt.Errorf("failed to connect: %w", err)
	}

	return driver, nil
}
