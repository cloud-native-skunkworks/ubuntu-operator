package main

import (
	"flag"
	"fmt"
	"net"
)

func echoServer(c net.Conn) {
	for {
		buf := make([]byte, 512)
		nr, err := c.Read(buf)
		if err != nil {
			return
		}

		data := buf[0:nr]
		fmt.Printf("Received: %v", string(data))
		_, err = c.Write(data)
		if err != nil {
			panic("Write: " + err.Error())
		}
	}
}

var (
	socketPath string
)

func main() {

	flag.StringVar(&socketPath, "socketPath", "", "socketPath")
	flag.Parse()

	fmt.Printf("Using socketpath %s", socketPath)
	l, err := net.Listen("unix", socketPath)
	if err != nil {
		println("listen error", err.Error())
		return
	}

	for {
		fd, err := l.Accept()
		if err != nil {
			println("accept error", err.Error())
			return
		}

		go echoServer(fd)
	}
}
