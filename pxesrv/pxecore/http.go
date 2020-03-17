package pxecore

import (
	"net"
	"net/http"
	"sync"
)

// HTTPStart start file server by http
func HTTPStart(wg *sync.WaitGroup) {
	defer func() {
		if err := recover(); err != nil {
			panic(err)
		}
	}()
	defer wg.Done()
	c := GetConf()
	listen := net.JoinHostPort(c.HTTP.HTTPIP, c.HTTP.HTTPPort)
	http.Handle("/", http.FileServer(http.Dir(c.HTTP.MountPath)))
	log.Printf("starting http server %s and handle on path: %s", listen, c.HTTP.MountPath)
	panic(http.ListenAndServe(listen, nil))
}
