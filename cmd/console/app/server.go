package console

import (
	"fmt"
	"github.com/spf13/cobra"
	utilerrors "k8s.io/apimachinery/pkg/util/errors"
	cliflag "k8s.io/component-base/cli/flag"
	"k8s.io/component-base/term"
	"v1/cmd/console/app/options"
	"v1/pkg/server"
	"v1/pkg/version/verflag"
)

func NewAPIServerCommand() (cmd *cobra.Command) {
	s := options.NewServerRunOptions()
	cmd = &cobra.Command{
		Use: "apiserver",
		Long: `The Tensor-v1-platform API server validates and configures data for the API objects. 
The API Server services REST operations and through which all other components interact.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			verflag.PrintAndExitIfRequested() // 不是报错（吧）

			if errs := s.Validate(); len(errs) != 0 {
				return utilerrors.NewAggregate(errs)
			}

			return Run(s, server.SetupSignalHandler())
		},
		SilenceUsage: true,
	}

	fs := cmd.Flags()
	namedFlagSets := s.Flags()
	verflag.AddFlags(namedFlagSets.FlagSet("global"))
	namedFlagSets.FlagSet("global").BoolP("help", "h", false, fmt.Sprintf("help for %s", cmd.Name()))
	for _, f := range namedFlagSets.FlagSets {
		fs.AddFlagSet(f)
	}

	cols, _, _ := term.TerminalSize(cmd.OutOrStdout())
	cliflag.SetUsageAndHelpFunc(cmd, namedFlagSets, cols)
	return
}

func Run(s *options.ServerRunOptions, stopCh <-chan struct{}) error {
	apiserver, err := s.NewAPIServer(stopCh)
	if err != nil {
		return err
	}

	err = apiserver.PrepareRun(stopCh)
	if err != nil {
		return nil
	}

	return apiserver.Run(stopCh)
}
