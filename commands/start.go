package commands

import (
	"context"
	"flag"
	"fmt"
	"time"

	"github.com/google/subcommands"
	"github.com/isksss/paperma-manager/config"
)

type StartCommand struct {
}

func (c *StartCommand) Name() string { return "start" }

func (c *StartCommand) Synopsis() string { return "start server." }

func (c *StartCommand) Usage() string { return "start" }

func (c *StartCommand) SetFlags(f *flag.FlagSet) {
}

func (c *StartCommand) Execute(ctx context.Context, f *flag.FlagSet, args ...interface{}) subcommands.ExitStatus {
	// Todo:
	// - javaの存在確認
	// - サブプロセスでPaperMC実行
	// - configで指定された時間に再起動

	// Get restart time.
	data, err := config.GetConfig()
	if err != nil {
		fmt.Printf("config read error: %v\n", err)
		return subcommands.ExitFailure
	}

	var parsedTime []time.Time
	for _, restartTime := range data.Server.RestartTime {
		parsed, err := time.Parse("15:04", restartTime)
		if err != nil {
			continue
		}

		parsedTime = append(parsedTime, parsed)
	}

	fmt.Println(parsedTime)
	return subcommands.ExitSuccess
}
