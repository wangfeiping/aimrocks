package commands

import (
	"github.com/spf13/cobra"
	"github.com/wangfeiping/aimrocks/log"
)

// NewTxCommand returns tx command
func NewTxCommand(run Runner) *cobra.Command {
	cmd := &cobra.Command{
		Use:   CmdTx,
		Short: "Transactions commands",
	}

	send := &cobra.Command{
		Use:   CmdTxSend,
		Short: "Send transaction(s) to ...",
		Run: func(cmd *cobra.Command, args []string) {
			if _, err := run(); err != nil {
				log.Error("Send tx error: ", err)
			}
			return
		},
	}

	send.Flags().String(FlagFrom, "", "One or more transfer out addresses")
	send.Flags().String(FlagFromAmount, "", "Amount of coins to transfer out")
	send.Flags().String(FlagTo, "", "One or more transfer in addresses")
	send.Flags().String(FlagToAmount, "", "Amount of coins to transfer in")
	send.Flags().Bool(FlagRelay, false,
		"Relay mode, transaction will be registered to the issuing chain")
	send.Flags().Bool(FlagTrustNode, false, "Don't verify proofs for responses")
	send.Flags().Int64(FlagMaxGas, 10000, "Max gas for transaction")
	cmd.AddCommand(send)

	cmd.SetUsageTemplate(usageTemplate)

	return cmd
}
