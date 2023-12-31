package commands

import (
	"context"
	"encoding/json"
	"flag"
	"io/ioutil"

	"github.com/google/subcommands"
	"github.com/isksss/paperm/config"
	"github.com/isksss/paperm/model"
)

type InitCommand struct {
	mode    string
	version string
	memory  int
}

func (c *InitCommand) Name() string { return "init" }

func (c *InitCommand) Synopsis() string { return "init config file." }

func (c *InitCommand) Usage() string { return "init [option]" }

func (c *InitCommand) SetFlags(f *flag.FlagSet) {
	f.StringVar(&c.mode, "mode", "paper", "server mode")
	f.StringVar(&c.version, "version", "1.20.1", "papermc version")
	f.IntVar(&c.memory, "memory", 1024, "server memory")
}

func (c *InitCommand) Execute(ctx context.Context, f *flag.FlagSet, args ...interface{}) subcommands.ExitStatus {
	configFile := model.Config{}
	configFile.Url = "https://papermc.io/api/v2/projects/"
	configFile.Mode = c.mode
	configFile.Server.MaxMemory = c.memory
	configFile.Server.MinMemory = c.memory
	configFile.PaperVersion = c.version
	configFile.WaterfallVersion = "1.20"
	configFile.Server.RestartTime = append(configFile.Server.RestartTime, "6:00")
	configFile.Server.RestartTime = append(configFile.Server.RestartTime, "18:00")
	configFile.JarName = "server.jar"
	configFile.Plugin.Download = true

	//plugin
	configFile.Plugin.Plugins = append(configFile.Plugin.Plugins,
		model.Plugins{
			Name: "geysermc.jar",
			Url:  "https://download.geysermc.org/v2/projects/geyser/versions/latest/builds/latest/downloads/spigot",
		},
		model.Plugins{
			Name: "floodgate.jar",
			Url:  "https://download.geysermc.org/v2/projects/floodgate/versions/latest/builds/latest/downloads/spigot",
		},
	)

	data, err := json.MarshalIndent(configFile, "", "  ")
	if err != nil {
		return subcommands.ExitFailure
	}

	err = ioutil.WriteFile(config.ConfigFileName, data, 0644)
	if err != nil {
		return subcommands.ExitFailure
	}

	// eula.txt
	eula := "eula=true"
	err = ioutil.WriteFile("eula.txt", []byte(eula), 0644)
	if err != nil {
		return subcommands.ExitFailure
	}

	return subcommands.ExitSuccess
}
