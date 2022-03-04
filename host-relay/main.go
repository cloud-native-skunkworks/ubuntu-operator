package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net"

	kmk "github.com/ElyKar/golang-kmod/kmod"
	"github.com/pmorjan/kmod"
)

type Module struct {
	Name  string `json:"name"`
	Flags string `json:"flags"`
}

type RelayMessage struct {
	Type           string   `json:"type"` // "Request | Response"
	DesiredModules []Module `json:"desiredModules"`
	ActualModules  []Module `json:"actualModules"`
}

func loadKernelModule(moduleName string, flags string, k *kmod.Kmod) error {

	return k.Load(moduleName, flags, 0)
}
func server(c net.Conn, k *kmod.Kmod, kkm *kmk.Kmod) {
	for {
		buf := make([]byte, 12000)
		nr, err := c.Read(buf)
		if err != nil {
			return
		}

		data := buf[0:nr]
		fmt.Printf("Received: %v", string(data))

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
			msg.ActualModules = actualModules
		}

		b, err := json.Marshal(msg)
		if err != nil {
			fmt.Println("Error:", err)
		}

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
		fmt.Printf("No --socketPath set")
		return
	}
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

		go server(fd, k, kkm)
	}
}
