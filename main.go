package main

import (
	"context"
	"flag"
	"os"

	"github.com/google/subcommands"
	"github.com/isksss/paperma-manager/commands"
)

func main() {
	run()
}

func init() {
	subcommands.Register(subcommands.HelpCommand(), "")
	subcommands.Register(subcommands.FlagsCommand(), "")
	subcommands.Register(subcommands.CommandsCommand(), "")
	subcommands.Register(&commands.InitCommand{}, "")
	subcommands.Register(&commands.DownloadCommand{}, "")

	flag.Parse()
}

func run() {
	ctx := context.Background()
	os.Exit(int(subcommands.Execute(ctx)))
}
