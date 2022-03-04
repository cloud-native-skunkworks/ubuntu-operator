package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"net"
	"os"

	kmk "github.com/ElyKar/golang-kmod/kmod"
	"github.com/pmorjan/kmod"
	log "github.com/sirupsen/logrus"
)

type Module struct {
	Name  string `json:"name"`
	Flags string `json:"flags"`
}

type RelayMessage struct {
	Type           string   `json:"type"` // "Request | Response"
	HostName       string   `json:"hostName"`
	DesiredModules []Module `json:"desiredModules"`
	ActualModules  []Module `json:"actualModules"`
}

func loadKernelModule(moduleName string, flags string, k *kmod.Kmod) error {

	return k.Load(moduleName, flags, 0)
}
func server(c net.Conn, k *kmod.Kmod, kkm *kmk.Kmod) {
	for {
		reader := bufio.NewReader(c)
		data, err := reader.ReadBytes('\n')
		if err != nil {
			println(err.Error())

			return
		}

		data = data[:len(data)-1]

		var msg RelayMessage
		err = json.Unmarshal(data, &msg)
		if err != nil {
			fmt.Println("Error unmarshalling message:", err)
			return
		}

		switch msg.Type {
		case "Request":

			for _, module := range msg.DesiredModules {
				fmt.Printf("\nDesired module: %s", module.Name)
			}

			for _, module := range msg.DesiredModules {
				if err := loadKernelModule(module.Name, module.Flags, k); err != nil {
					fmt.Println("Error loading module:", err)
					continue
				}
			}

			// List all loaded modules
			list, err := kkm.List()
			if err != nil {
				panic(err)
			}
			var actualModules []Module

			for _, module := range list {
				m := Module{Name: module.Name()}
				actualModules = append(actualModules, m)
			}
			hostname, err := os.Hostname()
			if err != nil {
				panic(err)
			}

			msg.Type = "Response"
			msg.HostName = hostname
			msg.ActualModules = actualModules
		}

		b, err := json.Marshal(msg)
		if err != nil {
			fmt.Println("Error:", err)
		}

		b = append(b, '\n')

		_, err = c.Write(b)
		if err != nil {
			panic("Write: " + err.Error())
		}
	}
}

var (
	socketPath string
)

func main() {

	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stdout)
	// Only log the warning severity or above.
	log.SetLevel(log.DebugLevel)

	k, err := kmod.New()
	if err != nil {
		log.Fatal(err)
	}
	kkm, err := kmk.NewKmod()

	if err != nil {
		panic(err)
	}

	flag.StringVar(&socketPath, "socketPath", "", "socketPath")
	flag.Parse()

	if socketPath == "" {
		log.Printf("No --socketPath set")
		return
	}
	log.Printf("Using socketpath %s", socketPath)
	l, err := net.Listen("unix", socketPath)
	if err != nil {
		log.Printf("listen error", err.Error())
		return
	}

	for {
		fd, err := l.Accept()
		if err != nil {
			log.Printf("accept error", err.Error())
			return
		}

		go server(fd, k, kkm)
	}
}
