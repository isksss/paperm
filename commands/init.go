package commands

import (
	"context"
	"encoding/json"
	"flag"
	"io/ioutil"

	"github.com/google/subcommands"
	"github.com/isksss/paperma-manager/config"
	"github.com/isksss/paperma-manager/model"
)

type InitCommand struct {
	version string
	memory  int
}

// これがサブコマンド名になる
func (c *InitCommand) Name() string { return "init" }

// コマンド一覧で出てくるサブコマンドの説明
func (c *InitCommand) Synopsis() string { return "init config file." }

// helpとかで出てくる使い方
func (c *InitCommand) Usage() string { return "init [option]" }

// flagライブラリでオプションの処理をするやつ
func (c *InitCommand) SetFlags(f *flag.FlagSet) {
	f.StringVar(&c.version, "version", "1.20.1", "papermc version")
	f.IntVar(&c.memory, "memory", 1024, "server memory")
}

// 本体
func (c *InitCommand) Execute(ctx context.Context, f *flag.FlagSet, args ...interface{}) subcommands.ExitStatus {
	configFile := model.Config{}
	configFile.Server.MaxMemory = c.memory
	configFile.Server.MinMemory = c.memory
	configFile.PaperVersion = c.version

	data, err := json.MarshalIndent(configFile, "", "  ")
	if err != nil {
		return subcommands.ExitFailure
	}

	err = ioutil.WriteFile(config.ConfigFileName, data, 0644)
	if err != nil {
		return subcommands.ExitFailure
	}

	return subcommands.ExitSuccess
}