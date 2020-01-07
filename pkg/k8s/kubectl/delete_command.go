package kubectl

import (
	"k8s.io/kubectl/pkg/cmd/delete"
	"strconv"
)

type deleteCommandFactory struct{}

var _ CommandFactory = &deleteCommandFactory{}

func (*deleteCommandFactory) create(config *CommandConfig) (commands []*Command, err error) {
	factory := createKubectlFactory(config)
	command := delete.NewCmdDelete(factory, config.Streams)
	setManifestPath(command, config.Params)
	setFlag(command, "timeout", strconv.Itoa(int(config.Params.CommandTimeout))+"s")

	commands = append(commands, &Command{command: command, args: []string{}})
	return
}
