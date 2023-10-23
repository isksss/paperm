package download

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"

	"github.com/isksss/paperm/config"
	"github.com/isksss/paperm/model"
)

func ServerDownload() error {
	configFile, err := config.GetConfig()
	if err != nil {
		return err
	}

	// select mode
	mode := configFile.Mode

	// version
	var version string
	if mode == "paper" {
		version = configFile.PaperVersion
	} else if mode == "waterfall" {
		version = configFile.WaterfallVersion
	} else {
		return fmt.Errorf("invalid mode: %s", mode)
	}

	// version check
	url := configFile.Url + mode + "/versions/" + version

	body, err := getApi(url)
	if err != nil {
		return fmt.Errorf("get api error: %v", err)
	}

	var project model.Project
	if err := json.Unmarshal(body, &project); err != nil {
		return fmt.Errorf("json unmarshal error: %v", err)
	}

	// if version is not found
	if project.Error != "" {
		return fmt.Errorf("version not found: %s", version)
	}

	// get latest build
	latestBuilds := project.Builds[len(project.Builds)-1]

	// builds詳細取得
	buildUrl := url + "/builds/" + strconv.Itoa(latestBuilds)
	body, err = getApi(buildUrl)
	if err != nil {
		return fmt.Errorf("get api error: %v", err)
	}

	var builds model.BuildProject
	if err := json.Unmarshal(body, &builds); err != nil {
		return fmt.Errorf("json unmarshal error: %v", err)
	}

	jarName := builds.Downloads.Application.Name
	jarUrl := buildUrl + "/downloads/" + jarName

	fmt.Printf("Download: %s\n", jarName)
	file, err := os.Create(configFile.JarName)
	if err != nil {
		return fmt.Errorf("create file error: %v", err)
	}
	defer file.Close()

	resp, err := http.Get(jarUrl)
	if err != nil {
		return fmt.Errorf("http get error: %v", err)
	}
	defer resp.Body.Close()

	_, err = io.Copy(file, resp.Body)
	if err != nil {
		return fmt.Errorf("io copy error: %v", err)
	}
	fmt.Printf("Downloaded: %s\n", jarName)
	return nil
}

func getApi(url string) ([]byte, error) {
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
