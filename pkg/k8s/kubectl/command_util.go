package kubectl

import (
	"github.com/kudoh/concourse-k8s-resource/pkg/k8s"
	"github.com/kudoh/concourse-k8s-resource/pkg/models"
	"github.com/spf13/cobra"
	"k8s.io/client-go/discovery/cached/memory"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/clientcmd/api"
	cmdutil "k8s.io/kubectl/pkg/cmd/util"
	"log"
)

func createKubectlFactory(cc *CommandConfig) cmdutil.Factory {

	restConfig, err := cc.ClientConfig.ClientConfig()
	if err != nil {
		log.Fatalln(err)
	}

	// enforce given namespace
	rawConfig, err := cc.ClientConfig.RawConfig()
	if err != nil {
		log.Fatalln(err)
	}
	newConfig := clientcmd.NewDefaultClientConfig(rawConfig,
		&clientcmd.ConfigOverrides{Context: api.Context{Namespace: cc.Namespace}})

	discoveryClient := memory.NewMemCacheClient(cc.Discovery)
	getter := k8s.NewConcourseRESTClientGetter(restConfig, discoveryClient, newConfig)
	factory := cmdutil.NewFactory(getter)

	return factory
}

func setManifestPath(command *cobra.Command, params *models.OutParams) {
	if params.Kustomize {
		setFlagArray(command, "kustomize", params.Paths...)
	} else {
		setFlagArray(command, "filename", params.Paths...)
	}
}

func setFlag(command *cobra.Command, name string, value string) {
	if err := command.Flags().Set(name, value); err != nil {
		log.Fatalln("Flag Set Failure", err)
	}
}

func setFlagArray(command *cobra.Command, name string, values ...string) {
	for _, v := range values {
		setFlag(command, name, v)
	}
}
