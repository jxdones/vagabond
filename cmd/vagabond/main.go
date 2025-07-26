package main

import "github.com/jxdones/vagabond/cli"

func main() {
	vagabond := cli.NewCli()
	cli.RegisterCommands(vagabond)
	vagabond.Run()
}
