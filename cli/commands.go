package cli

import "github.com/jxdones/vagabond/cmd"

func RegisterCommands(cli *CLI) {
	cli.RegisterCommand(Command{"create", "Create migrations files", cmd.Create})
	cli.RegisterCommand(Command{"pack", "Apply pending migrations", cmd.PackMigration})
	cli.RegisterCommand(Command{"unpack", "Rollback 1 or n migrations", cmd.UnpackMigrations})
	cli.RegisterCommand(Command{"help", "Show this help message", func(_ []string) error {
		cli.ShowHelp()
		return nil 
	}})
	cli.RegisterCommand(Command{"version", "Show the Vagabond version", cmd.ShowVersion})
}
