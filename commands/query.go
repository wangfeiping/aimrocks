package commands

import (
	"github.com/spf13/cobra"
	"github.com/wangfeiping/aimrocks/log"
)

// NewQueryCommand returns query command
func NewQueryCommand(run Runner, isKeepRunning bool) *cobra.Command {
	cmd := &cobra.Command{
		Use:   CmdQuery,
		Short: "Query commands",
		RunE: func(cmd *cobra.Command, args []string) error {
			log.Warn("not implemented yet!")
			return nil
		},
	}
	return cmd
}
