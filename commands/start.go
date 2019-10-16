package commands

import (
	"github.com/spf13/cobra"
	"github.com/wangfeiping/aimrocks/log"
)

// NewStartCommand returns start command
func NewStartCommand(run Runner) *cobra.Command {
	cmd := &cobra.Command{
		Use:   CmdStart,
		Short: "start the node of blockchain",
		Run: func(cmd *cobra.Command, args []string) {
			log.Warn("not implemented yet!")
		},
	}
	return cmd
}
