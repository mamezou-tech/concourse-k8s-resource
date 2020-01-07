package k8s

import (
	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/restmapper"
	"k8s.io/client-go/tools/clientcmd"
)

type ConcourseRESTClientGetter struct {
	restConfig      *rest.Config
	discoveryClient discovery.CachedDiscoveryInterface
	clientConfig    clientcmd.ClientConfig
}

var _ genericclioptions.RESTClientGetter = &ConcourseRESTClientGetter{}

func (f *ConcourseRESTClientGetter) ToRESTConfig() (*rest.Config, error) {
	return f.restConfig, nil
}

// ToDiscoveryClient returns discovery client
func (f *ConcourseRESTClientGetter) ToDiscoveryClient() (discovery.CachedDiscoveryInterface, error) {
	return f.discoveryClient, nil
}

// ToRESTMapper returns a restmapper
func (f *ConcourseRESTClientGetter) ToRESTMapper() (meta.RESTMapper, error) {
	// from ConfigFlags
	discoveryClient, err := f.ToDiscoveryClient()
	if err != nil {
		return nil, err
	}

	mapper := restmapper.NewDeferredDiscoveryRESTMapper(discoveryClient)
	expander := restmapper.NewShortcutExpander(mapper, discoveryClient)
	return expander, nil
}

// ToRawKubeConfigLoader return kubeconfig loader as-is
func (f *ConcourseRESTClientGetter) ToRawKubeConfigLoader() clientcmd.ClientConfig {
	return f.clientConfig
}

func NewConcourseRESTClientGetter(restConfig *rest.Config, discoveryClient discovery.CachedDiscoveryInterface,
	clientConfig clientcmd.ClientConfig) genericclioptions.RESTClientGetter {

	return &ConcourseRESTClientGetter{
		restConfig:      restConfig,
		discoveryClient: discoveryClient,
		clientConfig:    clientConfig,
	}
}
