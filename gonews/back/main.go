package main

import (
	"flag"

	"github.com/davidddw/gopj/gonews/back/common"
)

var configFileName = flag.String("c", "config.ini", "config file path (default config.ini)")

func main() {
	flag.Parse()
	config, _ := common.InitConfig(*configFileName)
	common.InitEnv(config)
	common.StartServ(config)
}
