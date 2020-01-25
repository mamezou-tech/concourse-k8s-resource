package kubectl

import (
	"github.com/mamezou-tech/concourse-k8s-resource/pkg/models"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/client-go/discovery"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"log"
	"sync"
)

type CommandConfig struct {
	Clientset    kubernetes.Interface
	Discovery    discovery.DiscoveryInterface
	ClientConfig clientcmd.ClientConfig
	Streams      genericclioptions.IOStreams
	Namespace    string
	Resources    []models.WatchResource
	Params       *models.OutParams
}

type Command struct {
	*cobra.Command
	args []string
}

type CommandFactory interface {
	create(config *CommandConfig) (commands []*Command, err error)
}

func RunCommand(factory CommandFactory, config *CommandConfig) error {
	var wg sync.WaitGroup
	cmds, err := factory.create(config)
	if err != nil {
		return err
	}
	wg.Add(len(cmds))
	for _, cmd := range cmds {
		printFlagsAndArgs(cmd)
		go run(*cmd, &wg)
	}
	wg.Wait()
	return nil
}

func run(cmd Command, wg *sync.WaitGroup) {
	defer wg.Done()
	log.Printf("running %s command... %v", cmd.Name(), cmd.args)
	cmd.Run(cmd.Command, cmd.args)
}

func NewCommandFactory(params *models.OutParams) CommandFactory {
	switch {
	case params.Undo:
		return &undoCommandFactory{}
	case params.Delete:
		return &deleteCommandFactory{}
	default:
		return &applyCommandFactory{}
	}
}

func printFlagsAndArgs(command *Command) {
	log.Println("** kubectl flags **")
	command.Flags().VisitAll(func(flag *pflag.Flag) {
		log.Printf("- %s -> %s\n", flag.Name, flag.Value.String())
	})
	log.Println("** kubectl args **")
	for _, arg := range command.args {
		log.Printf("- %s\n", arg)
	}
}
