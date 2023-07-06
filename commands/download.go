package commands

import (
	"context"
	"flag"

	"github.com/google/subcommands"
)

type DownloadCommand struct {
}

// これがサブコマンド名になる
func (c *DownloadCommand) Name() string { return "download" }

// コマンド一覧で出てくるサブコマンドの説明
func (c *DownloadCommand) Synopsis() string { return "download the papermc" }

// helpとかで出てくる使い方
func (c *DownloadCommand) Usage() string { return "download" }

// flagライブラリでオプションの処理をするやつ
func (c *DownloadCommand) SetFlags(f *flag.FlagSet) {

}

// 本体
func (c *DownloadCommand) Execute(ctx context.Context, f *flag.FlagSet, args ...interface{}) subcommands.ExitStatus {
	return subcommands.ExitSuccess
}
