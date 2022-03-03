package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"path/filepath"
	"time"
)

var (
	socketPath string
)

func main() {

	flag.StringVar(&socketPath, "socketPath", "", "socketPath")
	flag.Parse()

	files, err := ioutil.ReadDir(filepath.Dir(socketPath))
	if err != nil {
		log.Fatal(err)
	}

	for _, f := range files {
		fmt.Printf("--> %s", f.Name())
	}

	fmt.Printf("Using socketpath %s", socketPath)
	c, err := net.Dial("unix", socketPath)
	if err != nil {
		panic(err.Error())
	}
	for {
		_, err := c.Write([]byte("hi\n"))
		if err != nil {
			println(err.Error())
		}
		time.Sleep(time.Second * 5)
	}
}
