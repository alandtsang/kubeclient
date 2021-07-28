package config

import (
	"flag"
	"log"
	"path/filepath"

	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

// DefaultKubeConfig get the kubeconfig file in the .kube directory under the home.
func DefaultKubeConfig() (config *rest.Config, err error) {
	var cfg *string
	if home := homedir.HomeDir(); home != "" {
		cfg = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		cfg = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}
	flag.Parse()

	config, err = clientcmd.BuildConfigFromFlags("", *cfg)
	if err != nil {
		log.Printf("DefaultKubeConfig build kube config from flags failed, %s", err)
	}

	return config, nil
}
