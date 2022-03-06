package main

import (
	"bufio"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"net"
	"os"
	"strings"
	"time"

	v1alpha1 "github.com/cloud-native-skunkworks/ubuntu-operator/api/v1alpha1"
	log "github.com/sirupsen/logrus"
	mv1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

type Module struct {
	Name  string `json:"name"`
	Flags string `json:"flags"`
}

type RelayMessage struct {
	Type           string   `json:"type"` // "Request | Response"
	HostName       string   `json:"hostname"`
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
	myFlags            arrayFlags
	minWatchTimeout    = 5 * time.Minute
	masterURL          = flag.String("masterURL", "", "masterURL")
	kubeconfig         = flag.String("kubeconfig", "", "kubeconfig")
	socketPath         = flag.String("socketPath", "", "socketPath")
	SchemeGroupVersion = v1alpha1.GroupVersion
)

func addKnownTypes(scheme *runtime.Scheme) error {
	scheme.AddKnownTypes(SchemeGroupVersion,
		&v1alpha1.UbuntuMachineConfigurationList{},
		&v1alpha1.UbuntuMachineConfiguration{})

	mv1.AddToGroupVersion(scheme, SchemeGroupVersion)
	return nil
}

func buildClient(config *rest.Config) *rest.RESTClient {

	scheme := runtime.NewScheme()
	SchemeBuilder := runtime.NewSchemeBuilder(addKnownTypes)
	if err := SchemeBuilder.AddToScheme(scheme); err != nil {
		panic(err)
	}
	crdConfig := *config
	crdConfig.GroupVersion = &SchemeGroupVersion
	crdConfig.APIPath = "/apis"
	crdConfig.NegotiatedSerializer = serializer.NewCodecFactory(scheme)
	crdConfig.UserAgent = rest.DefaultKubernetesUserAgent()

	client, err := rest.RESTClientFor(&crdConfig)
	if err != nil {
		panic(err)
	}

	return client
}

func fetchUbuntuMachineConfigurationCR(restClient *rest.RESTClient) (v1alpha1.UbuntuMachineConfigurationList, error) {
	result := v1alpha1.UbuntuMachineConfigurationList{}
	err := restClient.Get().Resource("UbuntuMachineConfigurations").Do(context.TODO()).Into(&result)

	return result, err
}

func main() {

	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stdout)
	// Only log the warning severity or above.
	log.SetLevel(log.DebugLevel)

	flag.Var(&myFlags, "module", "Module and args e.g. --module=nvme_core --module=rfcomm=foo")
	flag.Parse()
	// Setup KubeClient -----------------------------------------------------------------------------

	kubeCfg, err := clientcmd.BuildConfigFromFlags(*masterURL, *kubeconfig)
	if err != nil {
		log.Fatalf("Error building kubeconfig: %s", err.Error())
	}

	restClient := buildClient(kubeCfg)

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

	// This allows the daemonset to pass through module lists
	if os.Getenv("MODULE_LIST") != "" {
		envList := strings.Split(os.Getenv("MODULE_LIST"), ",")
		for _, module := range envList {
			myFlags = append(myFlags, module)
		}
	}

	// ------------------------------------------------------------------------------------------------
	if len(myFlags) == 0 {
		log.Printf("No modules specified. Exiting.\n")
		return
	}
	if *socketPath == "" {
		log.Printf("No --socketPath set")
		return
	}

	for {
		log.Printf("Using socketpath %s\n", socketPath)
		c, err := net.Dial("unix", *socketPath)
		if err != nil {
			panic(err.Error())
		}
		defer c.Close()
		// // Check that the CR exists before we start polling
		// li, err := fetchUbuntuKernelModuleCR(restClient)
		// if err != nil || len(li.Items) == 0 {
		// 	log.Warningf("No UbuntuKernelModule CR found. Waiting for it to be created.")
		// 	continue
		// }
		// TODO:
		//
		// Currently the architecture is for a single UbuntuKernelModule CR.
		// This may need to change in the future
		//
		sendMessage := RelayMessage{
			Type:           "Request",
			DesiredModules: desiredModules,
		}
		log.Println("Marshalling message")
		b, err := json.Marshal(sendMessage)
		if err != nil {
			fmt.Println("Error:", err)
			continue
		}
		b = append(b, '\n')
		log.Println("Writing message")
		_, err = c.Write(b)
		if err != nil {
			log.Printf(err.Error())
			continue
		}
		log.Println("Reading response")
		reader := bufio.NewReader(c)
		data, err := reader.ReadBytes('\n')
		if err != nil {
			log.Printf(err.Error())
			continue
		}

		data = data[:len(data)-1]

		log.Println("Unmarshalling response")
		var msg RelayMessage
		err = json.Unmarshal(data, &msg)
		if err != nil {
			log.Printf("Error unmarshalling message:", err)
			return
		}

		switch msg.Type {
		case "Response":
			log.Printf("Response recieved \n")

			// Write back the changes
			li, err := fetchUbuntuMachineConfigurationCR(restClient)
			if err != nil || len(li.Items) == 0 {
				log.Printf("No UbuntuMachineConfiguration CR found. Waiting for it to be created.")
				time.Sleep(time.Second * 30)
				continue
			}
			//TODO:
			// Only interacting with the first CR

			UbuntuMachineConfigurationCR := li.Items[0]

			var modules []v1alpha1.Module
			for _, mods := range msg.ActualModules {

				modules = append(modules, v1alpha1.Module{
					Name:  mods.Name,
					Flags: mods.Flags,
				})

			}

			UbuntuMachineConfigurationCR.Status.Nodes = append(UbuntuMachineConfigurationCR.Status.Nodes, v1alpha1.Node{
				Name:    msg.HostName,
				Modules: modules,
			})

			// Update resource version ----------------------------------------------------------------
			UbuntuMachineConfigurationCR.Annotations = map[string]string{}
			// ----------------------------------------------------------------------------------------
			log.Println("Updating UbuntuMachineConfiguration CR")
			result := v1alpha1.UbuntuMachineConfiguration{}
			err = restClient.
				Put().
				Namespace(UbuntuMachineConfigurationCR.ObjectMeta.Namespace).
				Name(UbuntuMachineConfigurationCR.ObjectMeta.Name).
				Resource("UbuntuMachineConfigurations").
				Body(&UbuntuMachineConfigurationCR).
				SubResource("status").
				Do(context.TODO()).
				Into(&result)

			if err != nil {
				log.Errorf("Error updating UbuntuMachineConfiguration CR: %s", err.Error())
				continue
			}

		}

		time.Sleep(time.Second * 30)
	}
}
