package commands

import (
	"github.com/spf13/cobra"
	"github.com/wangfeiping/aimrocks/log"
)

// NewQueryCommand returns query command
func NewQueryCommand(run Runner) *cobra.Command {
	cmd := &cobra.Command{
		Use:   CmdQuery,
		Short: "Query commands",
		Run: func(cmd *cobra.Command, args []string) {
			log.Warn("not implemented yet!")
		},
	}
	return cmd
}
