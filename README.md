# Vagabond

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

Manage your migrations with confidence. Vagabond ensures smooth evolutions, easy rollbacks, and automatically keeps your schema documented

## Features
**Apply Migrations**: Executes pending migration scripts in chronological order.
**Rollback Migrations**: Reverts the last applied migration or a specified number of migrations.


## Usage

```bash
$ vagabond create your_new_migration
$ vagabond pack --dsn="./your_database.db" # apply migrations in a sqlite database
$ vagabond unpack --dsn="./your_database.db" # rollback migrations in a sqlite database
```

## Contributing

Contributions are welcome! Please follow these steps:

* Fork the repository.
* Create a new branch for your feature or bug fix.
* Make your changes and ensure they are well-tested.

Submit a pull request with a clear description of your changes.

## Don't forget to star!  ⭐

If you find this project helpful or are actively using it, please consider giving it a **star**. Thank you!