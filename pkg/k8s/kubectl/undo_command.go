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

const toRevisionName = "to_revision"

type undoCommandFactory struct{}

var _ CommandFactory = &undoCommandFactory{}

func (*undoCommandFactory) create(config *CommandConfig) (commands []*Command, err error) {
	factory := createKubectlFactory(config)

	for _, resource := range config.Resources {
		command := rollout.NewCmdRolloutUndo(factory, config.Streams)
		var args []string
		switch {
		case k8s.IsDeployment(resource.Kind):
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
			setFlag(command, toRevisionName, strconv.FormatInt(rev-1, 10))
			args = []string{fmt.Sprintf("%s/%s", "deployment", resource.Name)}

		case k8s.IsStatefulSet(resource.Kind):
			sts, err := config.Clientset.AppsV1().StatefulSets(config.Namespace).Get(resource.Name, metav1.GetOptions{})
			if err != nil {
				return nil, err
			}
			rev, err := config.Clientset.AppsV1().ControllerRevisions(config.Namespace).Get(sts.Status.CurrentRevision, metav1.GetOptions{})
			if err != nil {
				return nil, err
			}
			if rev.Revision == 1 {
				log.Fatalf("%s has rev 1, so undo operation cannot execute.\n", resource.Name)
			}
			log.Printf("[%s]current rev:%d\n", resource.Name, rev)
			setFlag(command, toRevisionName, strconv.FormatInt(rev.Revision-1, 10))
			args = []string{fmt.Sprintf("%s/%s", "statefulset", resource.Name)}
		default:
			log.Fatalln("unsupported resource kind", resource.Kind)
		}
		commands = append(commands, &Command{command: command, args: args})
	}
	return
}
