package cli

import cmd "github.com/jxdones/vagabond/commands"

func RegisterCommands(cli *CLI) {
	cli.RegisterCommand(Command{"create", "<name>", "create migrations files", cmd.Create})
	cli.RegisterCommand(Command{"pack", "", "apply pending migrations", cmd.PackMigration})
	cli.RegisterCommand(Command{"unpack", "[n]", "rollback last n migrations (default 1)", cmd.UnpackMigrations})
	cli.RegisterCommand(Command{"sketch", "[dir]", "dump the current database schema. (default dir: migrations)", cmd.SketchSchema})
	cli.RegisterCommand(Command{"help", "", "print this help message", func(_ []string) error {
		cli.ShowHelp()
		return nil
	}})
	cli.RegisterCommand(Command{"version", "", "print vagabond version", cmd.ShowVersion})
}
