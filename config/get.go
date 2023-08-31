package config

import (
	"encoding/json"
	"io/ioutil"

	"github.com/isksss/papermc-manager/model"
)

func GetConfig() (model.Config, error) {
	config := model.Config{}

	file, err := ioutil.ReadFile(ConfigFileName)
	if err != nil {
		return config, err
	}

	if err = json.Unmarshal(file, &config); err != nil {
		return config, err
	}

	return config, nil
}
