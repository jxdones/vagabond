package cli

import (
	"fmt"
	"os"
)

type Command struct {
	Name        string
	Usage       string
	Description string
	Execute     func(args []string) error
}

type CLI struct {
	commands    map[string]Command
	commandList []Command
}

func NewCli() *CLI {
	return &CLI{
		commands: make(map[string]Command),
	}
}

func (c *CLI) RegisterCommand(cmd Command) {
	c.commands[cmd.Name] = cmd
	c.commandList = append(c.commandList, cmd) // ensure the commands are always in the added order
}

func (c *CLI) Run() {
	if len(os.Args) < 2 {
		c.ShowHelp()
		os.Exit(1)
	}

	command := os.Args[1]
	cmd, exists := c.commands[command]
	if !exists {
		fmt.Println("Unknown command:", command)
		c.ShowHelp()
		os.Exit(1)
	}

	if err := cmd.Execute(os.Args[2:]); err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}

func (c *CLI) ShowHelp() {
	fmt.Println("Vagabond - Manage your migrations with confidence.")
	fmt.Println("Usage:")
	fmt.Println("  vagabond <command> [options]")
	fmt.Println("Commands:")
	for _, cmd := range c.commandList {
		argHint := ""
		if cmd.Usage != "" {
			argHint = " " + cmd.Usage
		}
		fmt.Printf("  %-15s %s\n", cmd.Name+argHint, cmd.Description)
	}

	fmt.Println("\nOptions:")
	fmt.Println("  --dsn      Database connection string (required)")
}
