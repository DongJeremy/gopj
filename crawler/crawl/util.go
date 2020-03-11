package crawl

import (
	"os"

	"github.com/Unknwon/goconfig"
)

func getConfig(filename string) (map[string]string, error) {
	cfg, err := goconfig.LoadConfigFile(filename)

	if err != nil {
		return nil, err
	}
	return cfg.GetSection("default")
}

func createDirIfNotExists(path string) error {
	if !isExist(path) {
		err := os.Mkdir(path, os.ModePerm)
		if err != nil {
			return err
		}
	}
	return nil
}

func isExist(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return false
}
