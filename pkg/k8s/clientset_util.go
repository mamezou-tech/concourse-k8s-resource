package k8s

import (
	"github.com/kudoh/concourse-k8s-resource/pkg/models"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/clientcmd/api"
	"log"
)

func NewClientSet(source *models.Source) (kubernetes.Interface, clientcmd.ClientConfig) {

	config := NewClientConfig(source)
	restConfig, err := config.ClientConfig()

	if err != nil {
		log.Fatal("cannot get rest client config:", err)
	}

	clientset, err := kubernetes.NewForConfig(restConfig)
	if err != nil {
		log.Fatal("cannot get k8s clientset:", err)
	}

	return clientset, config
}

func NewClientConfig(source *models.Source) clientcmd.ClientConfig {

	var config clientcmd.ClientConfig
	if source.ClientKey != "" && source.ClientCert != "" {
		config = createConfigFromClientCert(source)
	} else if source.ClientToken != "" {
		config = createConfigFromToken(source)
	} else if source.Kubeconfig != "" {
		config = createConfigFromKubeConfig(source)
	} else {
		log.Fatalln("unknown k8s client auth. check your source config.")
	}
	return config
}

func createConfigFromClientCert(source *models.Source) clientcmd.ClientConfig {
	log.Println("using client certificate")
	config := clientcmd.NewNonInteractiveClientConfig(
		api.Config{
			Clusters: map[string]*api.Cluster{"default": createCluster(source)},
			AuthInfos: map[string]*api.AuthInfo{"concourse": {
				ClientCertificateData: []byte(source.ClientCert),
				ClientKeyData:         []byte(source.ClientKey),
			}},
			Contexts: map[string]*api.Context{"default": createContext(source.Namespace)},
		},
		"default",
		&clientcmd.ConfigOverrides{},
		&clientcmd.ClientConfigLoadingRules{},
	)
	return config
}

func createConfigFromToken(source *models.Source) clientcmd.ClientConfig {
	log.Println("using client token")
	config := clientcmd.NewNonInteractiveClientConfig(
		api.Config{
			Clusters: map[string]*api.Cluster{"default": createCluster(source)},
			AuthInfos: map[string]*api.AuthInfo{"concourse": {
				Token: source.ClientToken,
			}},
			Contexts: map[string]*api.Context{"default": createContext(source.Namespace)},
		},
		"default",
		&clientcmd.ConfigOverrides{},
		&clientcmd.ClientConfigLoadingRules{},
	)
	return config
}

func createConfigFromKubeConfig(source *models.Source) clientcmd.ClientConfig {
	log.Println("using kubeconfig")
	config, err := clientcmd.Load([]byte(source.Kubeconfig))
	if err != nil {
		log.Fatalln("cannot create config from kubeconfig:", err)
	}
	clientConfig := clientcmd.NewDefaultClientConfig(*config, &clientcmd.ConfigOverrides{})
	return clientConfig
}

func createCluster(source *models.Source) *api.Cluster {
	return &api.Cluster{
		Server:                   source.ApiServerUrl,
		CertificateAuthorityData: []byte(source.ApiServerCA),
		InsecureSkipTLSVerify:    source.SkipTLSVerify,
	}
}

func createContext(ns string) *api.Context {
	return &api.Context{
		AuthInfo:  "concourse",
		Cluster:   "default",
		Namespace: ns,
	}
}
