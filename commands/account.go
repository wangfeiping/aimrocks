package commands

import (
	"github.com/spf13/cobra"
	"github.com/wangfeiping/aimrocks/log"
)

// NewAccountCommand returns account command
func NewAccountCommand(run Runner) *cobra.Command {
	cmd := &cobra.Command{
		Use:   CmdAccount,
		Short: "Account commands",
		Run: func(cmd *cobra.Command, args []string) {
			log.Warn("not implemented yet!")
		},
	}

	cmd.SetUsageTemplate(usageTemplate)

	return cmd
}
