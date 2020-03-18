package main

import (
	"log"

	"github.com/davidddw/gopj/pxesrv/pxecore"
)

func main() {
	log.Printf("starting pxe server...")
	serve := pxecore.Server{Config: pxecore.GetConf("pxe.yml")}
	serve.Serve()
}
