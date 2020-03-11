package main

import (
	"log"

	"github.com/Unknwon/goconfig"
)

// Config get conf from file
type Config struct {
	URL       string
	Chapter   string
	PhpURL    string
	URLFile   string
	NovelFile string
	StartPage int
}

var config Config

func init() {
	config = getConfig("app.conf")
}

func getConfig(filename string) (c Config) {
	cfg, err := goconfig.LoadConfigFile(filename)
	if err != nil {
		log.Fatal(err)
	}
	baseURL := cfg.MustValue("", "baseurl", "")
	chapterIndex := cfg.MustValue("", "chapterindex", "")
	phpRequest := cfg.MustValue("", "phprequest", "")
	startPage := cfg.MustInt("", "startpage", 2001)
	urlFile := cfg.MustValue("", "urlfile", "1.txt")
	novelFile := cfg.MustValue("", "novelfile", "2.txt")
	c = Config{
		baseURL,
		chapterIndex,
		phpRequest,
		urlFile,
		novelFile,
		startPage,
	}
	return
}
