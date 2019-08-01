package commands

import (
	"github.com/spf13/cobra"
	"github.com/wangfeiping/aimrocks/log"
)

// NewTxCommand returns tx command
func NewTxCommand(run Runner, isKeepRunning bool) *cobra.Command {
	cmd := &cobra.Command{
		Use:   CmdTx,
		Short: "Transactions commands",
		RunE: func(cmd *cobra.Command, args []string) error {
			log.Warn("not implemented yet!")
			return nil
		},
	}
	return cmd
}
