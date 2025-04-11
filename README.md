# Vagabond

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

Manage your migrations with confidence. Vagabond ensures smooth evolutions, easy rollbacks, and automatically keeps your schema documented

## Features
* Apply Migrations: Executes pending migration scripts in chronological order.
* Rollback Migrations: Reverts the last applied migration or a specified number of migrations.
* Migration Tracking: Keeps track of applied migrations in a dedicated database table.
* Schema Dumping: Generates a schema.sql file reflecting the current database schema.

## Usage
Vagabond currently supports SQLite only. Support for MariaDB/MySQL and PostgreSQL will be added soon.

```bash
$ vagabond
Vagabond - Manage your migrations with confidence.
Usage:
  vagabond <command> [options]
Commands:
  create <name>   create migrations files
  pack            apply pending migrations
  unpack [n]      rollback last n migrations (default 1)
  sketch [dir]    dump the current database schema. (default dir: migrations)
  help            print this help message
  version         print vagabond version

Options:
  --dsn      Database connection string (required)
$ vagabond create your_new_migration
$ vagabond pack --dsn="./your_database.db"
$ vagabond unpack --dsn="./your_database.db"
```

## Contributing

Contributions are welcome! Please follow these steps:

* Fork the repository.
* Create a new branch for your feature or bug fix.
* Make your changes and ensure they are well-tested.

Submit a pull request with a clear description of your changes.

## Don't forget to star!  ⭐

If you find this project helpful or are actively using it, please consider giving it a **star**. Thank you!