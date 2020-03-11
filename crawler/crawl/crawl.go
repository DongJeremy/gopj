package crawl

import (
	"fmt"
	"log"
	"os"
	"runtime"
)

const configFile = "conf/app.conf"

var (
	ch1      chan string
	ch2      chan string
	ch3      chan int
	novelDir string
	config   map[string]string
)

func init() {
	ch1 = make(chan string, 20)
	ch2 = make(chan string, 1000)
	ch3 = make(chan int, 1000)

	config, _ := getConfig(configFile)

	logfile, err := os.OpenFile(config["LOG_FILE"], os.O_CREATE|os.O_RDWR|os.O_APPEND, 0777)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	log.SetOutput(logfile)
}

// DownloadNovels export
func DownloadNovels() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	createDirIfNotExists(config["DOWNLOAD_FOLDER"])

	go getChapterListURL(config["NOVEL_URL"])
	go getChapterNumber(config["SUFFIX"])
	go downloadNovel()

	count := 0
	for num := range ch3 {
		count = count + num
		fmt.Println("count:", count)
	}
	fmt.Println("crawl end")
}
