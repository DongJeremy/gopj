package common

import (
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/sirupsen/logrus"
)

// News news structure
type News struct {
	ID    int64
	Title string
	Link  url.URL
	Ctime time.Time
}

var (
	logger *logrus.Logger
)

// InitEnv init env variable
func InitEnv(config *Config) {
	initRedis(config)
	logger = LoggerFromConfig(config)
}

func execute(workpath, path string, args ...string) ([]byte, error) {
	cmd := exec.Command(path, args...)
	cmd.Dir = workpath
	return cmd.Output()
}

// InitDataPuller load data from file and save into redis
// dir E:/tmp/data
// repo https://github.com/gocn/news.git
func InitDataPuller(config *Config) error {
	dir := config.Common.DataFolder
	err := os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		return err
	}
	folderPath, _ := filepath.Abs(dir + "/news")
	if _, err := os.Stat(folderPath); os.IsNotExist(err) {
		out, err1 := execute(dir, "git", "clone", config.Common.Repo)
		if err1 != nil {
			logger.Errorf("%s", "clone failed")
			logger.Errorf("%s", err)
			logger.Errorf("%s", out)
			return err1
		}
		logger.Infof("%s", "Success to clone news")
	}
	out, err := execute(folderPath, "git", "pull", "origin", "master")
	if err != nil {
		logger.Errorf("%s", "Pull failed")
		logger.Errorf("%s", err)
		logger.Errorf("%s", out)
		return err
	}
	logger.Infof("%s", "Success to pull news")

	// 缓存数据操作
	files := GetFileList(folderPath)
	for _, file := range files {
		wg.Add(1)
		go CacheNews(file)
	}
	wg.Wait()
	logger.Infof("%s", "Success to cache news")

	return nil
}
