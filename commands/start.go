package commands

import (
	"github.com/spf13/cobra"
	"github.com/wangfeiping/aimrocks/log"
)

// NewStartCommand returns start command
func NewStartCommand(run Runner, isKeepRunning bool) *cobra.Command {
	cmd := &cobra.Command{
		Use:   CmdStart,
		Short: "Start the node",
		RunE: func(cmd *cobra.Command, args []string) error {
			log.Warn("not implemented yet!")
			return nil
		},
	}
	return cmd
}
