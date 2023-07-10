package commands

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"

	"github.com/google/subcommands"
	"github.com/isksss/paperma-manager/commands/download"
	"github.com/isksss/paperma-manager/config"
	"github.com/isksss/paperma-manager/model"
)

type DownloadCommand struct {
}

func (c *DownloadCommand) Name() string { return "download" }

func (c *DownloadCommand) Synopsis() string { return "download the papermc" }

func (c *DownloadCommand) Usage() string { return "download" }

func (c *DownloadCommand) SetFlags(f *flag.FlagSet) {
}

func (c *DownloadCommand) Execute(ctx context.Context, f *flag.FlagSet, args ...interface{}) subcommands.ExitStatus {

	// get config data
	data, err := config.GetConfig()
	if err != nil {
		fmt.Printf("config read error: %v\n", err)
		return subcommands.ExitFailure
	}

	// download papermc
	version := data.PaperVersion // Papermc version

	// バージョン確認
	// 指定したバージョンが有効か検証
	url := config.ApiUrl + "/versions/" + version

	body, err := GetApi(url)
	if err != nil {
		return subcommands.ExitFailure
	}

	var project model.Project
	if err := json.Unmarshal(body, &project); err != nil {
		return subcommands.ExitFailure
	}

	// 指定されたバージョン存在しなければreturn
	if project.Error != "" {
		return subcommands.ExitFailure
	}

	// 最新のビルドを取得
	latestBuilds := project.Builds[len(project.Builds)-1]

	// builds詳細取得
	buildUrl := url + "/builds/" + strconv.Itoa(latestBuilds)
	body, err = GetApi(buildUrl)
	if err != nil {
		return subcommands.ExitFailure
	}

	var builds model.BuildProject
	if err := json.Unmarshal(body, &builds); err != nil {
		return subcommands.ExitFailure
	}

	// jarfileのダウンロード
	jarName := builds.Downloads.Application.Name
	jarUrl := buildUrl + "/downloads/" + jarName

	// create file
	file, err := os.Create(data.JarName)
	if err != nil {
		return subcommands.ExitFailure
	}
	defer file.Close()

	resp, err := http.Get(jarUrl)
	if err != nil {
		return subcommands.ExitFailure
	}
	defer resp.Body.Close()

	_, err = io.Copy(file, resp.Body)
	if err != nil {
		return subcommands.ExitFailure
	}

	//Download
	err = download.Download()
	if err != nil {
		return subcommands.ExitFailure
	}

	return subcommands.ExitSuccess
}

func GetApi(url string) ([]byte, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}
