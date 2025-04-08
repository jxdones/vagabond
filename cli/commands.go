package cli

import "github.com/jxdones/vagabond/cmd"

func RegisterCommands(cli *CLI) {
	cli.RegisterCommand(Command{"create", "create migrations files", cmd.Create})
	cli.RegisterCommand(Command{"pack", "apply pending migrations", cmd.PackMigration})
	cli.RegisterCommand(Command{"unpack", "rollback 1 or n migrations", cmd.UnpackMigrations})
	cli.RegisterCommand(Command{"sketch", "dump the current database schema. If no directory is provided, it will be created under the migrations directory", cmd.SketchSchema},)
	cli.RegisterCommand(Command{"help", "print this help message", func(_ []string) error {
		cli.ShowHelp()
		return nil 
	}})
	cli.RegisterCommand(Command{"version", "print vagabond version", cmd.ShowVersion})
}
