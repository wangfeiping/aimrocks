package commands

import (
	"github.com/spf13/cobra"
	"github.com/wangfeiping/aimrocks/log"
)

// NewKeyCommand returns key command
func NewKeyCommand(run Runner) *cobra.Command {
	cmd := &cobra.Command{
		Use:   CmdKey,
		Short: "key commands",
		Run: func(cmd *cobra.Command, args []string) {
			log.Warn("not implemented yet!")
		},
	}

	cmd.SetUsageTemplate(usageTemplate)

	return cmd
}
