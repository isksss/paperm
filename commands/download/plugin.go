package download

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"sync"

	"github.com/isksss/papermc-manager/config"
)

func Download() error {
	// get config data
	data, err := config.GetConfig()
	if err != nil {
		return err
	}

	if !data.Plugin.Download {
		err := fmt.Errorf("plugin download is disabled")
		return err
	}

	if _, err := os.Stat("plugins"); os.IsNotExist(err) {
		if err := os.MkdirAll("plugins", 0755); err != nil {
			return err
		}
	}

	// download plugin
	var wg sync.WaitGroup
	for _, plugin := range data.Plugin.Plugins {
		// download plugin
		name := plugin.Name
		url := plugin.Url
		wg.Add(1)
		go downloadPlugin(name, url, &wg)
	}

	wg.Wait()
	fmt.Println("Downloaded all plugins")
	return nil
}

func downloadPlugin(name string, url string, wg *sync.WaitGroup) error {
	defer wg.Done()
	fmt.Printf("Downloading %s...\n", name)
	dir := filepath.Join("plugins", name)
	fileName, err := os.Create(dir)
	if err != nil {
		return err
	}

	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	_, err = io.Copy(fileName, resp.Body)
	if err != nil {
		return err
	}
	fmt.Printf("Downloaded %s\n", name)
	return nil
}
