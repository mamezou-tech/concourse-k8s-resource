package kubectl

import (
	"github.com/emicklei/go-restful/log"
	"github.com/spf13/cobra"
	"k8s.io/client-go/util/exec"
	"k8s.io/kubectl/pkg/cmd/apply"
	"k8s.io/kubectl/pkg/cmd/diff"
	cmdutil "k8s.io/kubectl/pkg/cmd/util"
)

type diffCommandFactory struct{}

var _ CommandFactory = &diffCommandFactory{}

func (*diffCommandFactory) create(config *CommandConfig) (commands []*Command, err error) {
	factory := createKubectlFactory(config)

	options := diff.NewDiffOptions(config.Streams)
	command := &cobra.Command{
		Use:                   "diff",
		DisableFlagsInUseLine: true,
		Run: func(cmd *cobra.Command, args []string) {
			cmdutil.CheckDiffErr(options.Complete(factory, cmd))
			if err := options.Run(); err != nil {
				// exit code == 1 -> there is difference(not error!)
				if ee, ok := err.(exec.ExitError); ok && ee.ExitStatus() == 1 {
					log.Printf("found difference!")
					return
				}
				// exit with error code(>2)
				log.Printf("ERR! %+v", err)
				cmdutil.CheckDiffErr(err)
			}
		},
	}
	cmdutil.AddFilenameOptionFlags(command, &options.FilenameOptions, "contains the configuration to diff")
	cmdutil.AddServerSideApplyFlags(command)
	cmdutil.AddFieldManagerFlagVar(command, &options.FieldManager, apply.FieldManagerClientSideApply)

	setManifestPath(command, config.Params)

	commands = append(commands, &Command{command, []string{}})
	return
}
