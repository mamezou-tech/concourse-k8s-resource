package kubectl

import (
	"k8s.io/kubectl/pkg/cmd/apply"
	"strconv"
)

type applyCommandFactory struct{}

var _ CommandFactory = &applyCommandFactory{}

func (*applyCommandFactory) create(config *CommandConfig) (commands []*Command, err error) {
	factory := createKubectlFactory(config)
	command := apply.NewCmdApply("kubectl", factory, config.Streams)
	setFlag(command, "record", "true")
	setFlag(command, "timeout", strconv.Itoa(int(config.Params.CommandTimeout))+"s")
	setManifestPath(command, config.Params)

	commands = append(commands, &Command{command: command, args: []string{}})
	return
}
