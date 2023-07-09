package commands

import (
	"context"
	"flag"

	"github.com/google/subcommands"
)

type StartCommand struct {
}

func (c *StartCommand) Name() string { return "start" }

func (c *StartCommand) Synopsis() string { return "start server." }

func (c *StartCommand) Usage() string { return "start" }

func (c *StartCommand) SetFlags(f *flag.FlagSet) {
}

func (c *StartCommand) Execute(ctx context.Context, f *flag.FlagSet, args ...interface{}) subcommands.ExitStatus {
	return subcommands.ExitSuccess
}
