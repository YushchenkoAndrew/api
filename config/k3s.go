package config

import (
	"flag"
	"os"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

var K3s *kubernetes.Clientset

func LoadK3s() {
	var config *rest.Config
	var err error

	if os.Getenv("GIN_MODE") == "release" {
		config, err = rest.InClusterConfig()
	} else {
		kubeConfig := flag.String("kubeconfig", "./k3s.yaml", "kubeconfig file location")
		config, err = clientcmd.BuildConfigFromFlags("", *kubeConfig)
	}

	if err != nil {
		panic(err.Error())
	}

	if K3s, err = kubernetes.NewForConfig(config); err != nil {
		panic(err.Error())
	}
}
