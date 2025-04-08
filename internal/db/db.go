package db

type Driver interface {
	Connect(dsn string) error
	Close() error
	GetAppliedMigrations() (map[string]bool, error)
	GetAppliedMigrationsList() ([]string, error)
	ExecuteMigration(filePath string) error
	RollbackMigration(filePath, name string) error
	DumpSchema() (string, error)
}
