package kubectl

import (
	"fmt"
	"github.com/kudoh/concourse-k8s-resource/pkg/k8s"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/kubectl/pkg/cmd/rollout"
	"k8s.io/kubectl/pkg/util/deployment"
	"log"
	"strconv"
)

type undoCommandFactory struct{}

var _ CommandFactory = &undoCommandFactory{}

func (*undoCommandFactory) create(config *CommandConfig) (commands []*Command, err error) {
	factory := createKubectlFactory(config)

	for _, resource := range config.Resources {
		switch {
		case k8s.IsDeployment(resource.Kind):
			command := rollout.NewCmdRolloutUndo(factory, config.Streams)
			d, err := config.Clientset.AppsV1().Deployments(config.Namespace).Get(resource.Name, metav1.GetOptions{})
			if err != nil {
				return nil, err
			}
			rev, err := deployment.Revision(d)
			if err != nil {
				return nil, err
			}
			if rev == 1 {
				log.Fatalf("%s has rev 1, so undo operation cannot execute.\n", resource.Name)
			}
			log.Printf("[%s]current rev:%d\n", resource.Name, rev)
			setFlag(command, "to-revision", strconv.Itoa(int(rev-1)))
			args := []string{fmt.Sprintf("%s/%s", "deployment", resource.Name)}

			commands = append(commands, &Command{command: command, args: args})
		case k8s.IsStatefulSet(resource.Kind):
			log.Printf("[Warning] Sorry, currently undo operation does not support StatefulSet[%s]. skip...\n", resource.Name)
			continue
		default:
			log.Fatalln("unsupported resource kind", resource.Kind)
		}
	}
	return
}
