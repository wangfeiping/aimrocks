package commands

import (
	"github.com/spf13/cobra"
	"github.com/wangfeiping/aimrocks/log"
)

// NewAccountCommand returns account command
func NewAccountCommand(run Runner, isKeepRunning bool) *cobra.Command {
	cmd := &cobra.Command{
		Use:   CmdAccount,
		Short: "Account commands",
		RunE: func(cmd *cobra.Command, args []string) error {
			log.Warn("not implemented yet!")
			return nil
		},
	}
	return cmd
}
