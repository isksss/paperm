package commands

import (
	"context"
	"flag"
	"fmt"
	"os/exec"
	"strconv"
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

	// exists java
	javaCmd := "java"
	javaBin, err := exec.LookPath(javaCmd)
	if err != nil {
		return subcommands.ExitFailure
	}

	//papermcの起動
	// opt
	min := strconv.Itoa(data.Server.MinMemory)
	max := strconv.Itoa(data.Server.MaxMemory)
	xms := "-Xms" + min + "M"
	xmx := "-Xmx" + max + "M"
	jar := data.JarName

	// server
	cmd := exec.Command(javaBin, "-jar", xms, xmx, jar)
	err = cmd.Start()
	if err != nil {
		return subcommands.ExitFailure
	}

	cmd.Wait()
	//debug
	fmt.Println("javabin", javaBin)
	fmt.Println(parsedTime)
	return subcommands.ExitSuccess
}
