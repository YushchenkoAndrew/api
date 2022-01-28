package config

import (
	"flag"
	"os"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	metrics "k8s.io/metrics/pkg/client/clientset/versioned"
)

var K3s *kubernetes.Clientset
var K3sConfig *rest.Config
var Metrics *metrics.Clientset

func LoadK3s() {
	// var config *rest.Config
	var err error

	if os.Getenv("GIN_MODE") == "release" {
		K3sConfig, err = rest.InClusterConfig()
	} else {
		kubeConfig := flag.String("kubeconfig", "./k3s.yaml", "kubeconfig file location")
		K3sConfig, err = clientcmd.BuildConfigFromFlags("", *kubeConfig)
	}

	if err != nil {
		panic(err.Error())
	}

	if K3s, err = kubernetes.NewForConfig(K3sConfig); err != nil {
		panic(err.Error())
	}

	if Metrics, err = metrics.NewForConfig(K3sConfig); err != nil {
		panic(err.Error())
	}
}
