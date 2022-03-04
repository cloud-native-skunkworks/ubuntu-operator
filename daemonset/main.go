package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"math/rand"
	"net"
	"strings"
	"time"

	ukmv1alpha1 "github.com/cloud-native-skunkworks/generated/ubuntukernelmodules/clientset/versioned"
	log "github.com/sirupsen/logrus"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
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

type arrayFlags []string

func (i *arrayFlags) String() string {
	return "my string representation"
}

func (i *arrayFlags) Set(value string) error {
	*i = append(*i, value)
	return nil
}

var (
	myFlags         arrayFlags
	minWatchTimeout = 5 * time.Minute
	timeoutSeconds  = int64(minWatchTimeout.Seconds() * (rand.Float64() + 1.0))
	masterURL       string
	kubeconfig      string
	socketPath      = flag.String("socketPath", "", "socketPath")
	addr            = flag.String("listen-address", ":8080", "The address to listen on for HTTP requests.")
)

func main() {
	flag.Var(&myFlags, "module", "Module and args e.g. -module=nvme_core --module=rfcomm=foo")
	flag.Parse()
	// Setup KubeClient -----------------------------------------------------------------------------

	kubeCfg, err := clientcmd.BuildConfigFromFlags(masterURL, kubeconfig)
	if err != nil {
		log.Fatalf("Error building kubeconfig: %s", err.Error())
	}

	log.Info("Built config from flags...")

	_, err = kubernetes.NewForConfig(kubeCfg)
	if err != nil {
		log.Fatalf("Error building watcher clientset: %s", err.Error())
	}

	ukmv1alpha1.NewForConfigOrDie(kubeCfg)

	// Setup Kmod ------------------------------------------------------------------------------------

	var desiredModules []Module
	for _, module := range myFlags {
		fmt.Println("Desired module:", module)
		s := strings.Split(module, "=")
		m := Module{Name: s[0]}
		if len(s) > 1 {
			m.Flags = s[1]
		}
		desiredModules = append(desiredModules, m)
	}

	if len(myFlags) == 0 {
		fmt.Printf("No modules specified. Exiting.\n")
		return
	}
	if *socketPath == "" {
		fmt.Printf("No --socketPath set")
		return
	}
	fmt.Printf("Using socketpath %s", socketPath)
	c, err := net.Dial("unix", *socketPath)
	if err != nil {
		panic(err.Error())
	}
	for {

		sendMessage := RelayMessage{
			Type:           "Request",
			DesiredModules: desiredModules,
		}

		b, err := json.Marshal(sendMessage)
		if err != nil {
			fmt.Println("Error:", err)
			continue
		}

		_, err = c.Write(b)
		if err != nil {
			println(err.Error())
		}

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
		case "Response":
			fmt.Printf("Response: %s", msg.ActualModules)
		}

		time.Sleep(time.Second * 30)
	}
}
