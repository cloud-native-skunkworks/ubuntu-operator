package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"net"
	"os"
	"os/exec"

	kmk "github.com/ElyKar/golang-kmod/kmod"
	"github.com/fatih/color"
	"github.com/pmorjan/kmod"
	log "github.com/sirupsen/logrus"
)

type Module struct {
	Name  string `json:"name"`
	Flags string `json:"flags"`
}

type DesiredPackages struct {
	Apt  []AptPackage  `json:"apt"`
	Snap []SnapPackage `json:"snap"`
}

type AptPackage struct {
	Name string `json:"name"`
}

type SnapPackage struct {
	Name        string `json:"name"`
	Confinement string `json:"confinement"`
}

type RelayMessage struct {
	Type            string          `json:"type"` // "Request | Response"
	HostName        string          `json:"hostname"`
	DesiredModules  []Module        `json:"desiredModules"`
	DesiredPackages DesiredPackages `json:"desiredPackages"`
	ActualModules   []Module        `json:"actualModules"`
}

func loadKernelModule(moduleName string, flags string, k *kmod.Kmod) error {

	color.Blue(fmt.Sprintf("\nLoading kernel module: %s", moduleName))
	return k.Load(moduleName, flags, 0)
}

func loadAPTPackage(name string) error {
	var cmd *exec.Cmd
	var err error
	color.Blue(fmt.Sprintf("\nLoading Apt package: %s", name))
	cmd = exec.Command("apt", "install", name, "-y")
	stderr, _ := cmd.StderrPipe()
	if err = cmd.Start(); err != nil {
		log.Fatal(err)
	}

	scanner := bufio.NewScanner(stderr)
	for scanner.Scan() {
		fmt.Println(scanner.Text())
	}
	return err
}

func loadSnapPackage(name string, confinement string) error {
	var cmd *exec.Cmd
	var err error
	color.Blue(fmt.Sprintf("\nLoading Snap package: %s", name))
	if confinement == "classic" {
		cmd = exec.Command("snap", "install", name, "--classic")
	} else {
		cmd = exec.Command("snap", "install", name)
	}
	stderr, _ := cmd.StderrPipe()
	if err = cmd.Start(); err != nil {
		log.Fatal(err)
	}

	scanner := bufio.NewScanner(stderr)
	for scanner.Scan() {
		fmt.Println(scanner.Text())
	}
	return err
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
				if err := loadKernelModule(module.Name, module.Flags, k); err != nil {
					fmt.Println("Error loading module:", err)
					continue
				}
			}

			for _, pkg := range msg.DesiredPackages.Apt {
				fmt.Printf("\nDesired APT package: %s", pkg.Name)
				if err := loadAPTPackage(pkg.Name); err != nil {
					fmt.Println("Error loading package:", err)
					continue
				}
			}

			for _, pkg := range msg.DesiredPackages.Snap {
				fmt.Printf("\nDesired SNAP package: %s confinement type: %s", pkg.Name, pkg.Confinement)
				if err := loadSnapPackage(pkg.Name, pkg.Confinement); err != nil {
					fmt.Println("Error loading package:", err)
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
