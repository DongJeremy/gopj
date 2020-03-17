package pxecore

import (
	"fmt"
	"io"
	"os"
	"sync"
	"time"

	"github.com/pin/tftp"
)

// readHandler is called when client starts file download from server
func readHandler(filename string, rf io.ReaderFrom) error {
	c := GetConf()
	file, err := os.Open(c.TFTP.TftpPath + "/" + filename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		return err
	}
	n, err := rf.ReadFrom(file)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		return err
	}
	log.Printf("DHTP: tftp_files %d bytes sent\n", n)
	return nil
}

// writeHandler is called when client starts file upload to server
func writeHandler(filename string, wt io.WriterTo) error {
	c := GetConf()
	file, err := os.OpenFile(c.TFTP.TftpPath+"/"+filename, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0644)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		return err
	}
	n, err := wt.WriteTo(file)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		return err
	}
	log.Printf("DHTP: tftp_files %d bytes received\n", n)
	log.Printf("DHTP: tftp_files recieved and stored file to %s", c.TFTP.TftpPath+"/"+filename)
	return nil
}

// TFTPStart start server
func TFTPStart(wg *sync.WaitGroup) {
	defer func() {
		if err := recover(); err != nil {
			panic(err)
		}
	}()
	defer wg.Done()
	s := tftp.NewServer(readHandler, writeHandler)
	s.SetTimeout(5 * time.Second) // optional
	log.Printf("starting tftp_files server and listening on port :69 handle on path: %s", GetConf().TFTP.TftpPath)
	err := s.ListenAndServe("0.0.0.0:69") // blocks until s.Shutdown() is called
	if err != nil {
		fmt.Fprintf(os.Stdout, "server: %v\n", err)
		os.Exit(1)
	}
}
