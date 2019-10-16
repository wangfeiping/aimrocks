package commands

import (
	"github.com/spf13/cobra"
	"github.com/wangfeiping/aimrocks/log"
)

// NewTxCommand returns tx command
func NewTxCommand(run Runner) *cobra.Command {
	cmd := &cobra.Command{
		Use:   CmdTx,
		Short: "transactions commands",
	}

	send := &cobra.Command{
		Use:   CmdTxSend,
		Short: "send transaction(s) to ...",
		Run: func(cmd *cobra.Command, args []string) {
			if _, err := run(); err != nil {
				log.Error("Send tx error: ", err)
			}
			return
		},
	}

	send.Flags().String(FlagFrom, "", "one or more transfer out addresses")
	send.Flags().String(FlagFromAmount, "", "amount of coins to transfer out")
	send.Flags().String(FlagTo, "", "one or more transfer in addresses")
	send.Flags().String(FlagToAmount, "", "amount of coins to transfer in")
	send.Flags().Bool(FlagRelay, false,
		"relay mode, transaction will be registered to the issuing chain")
	send.Flags().Bool(FlagTrustNode, false, "don't verify proofs for responses")
	send.Flags().Int64(FlagMaxGas, 10000, "max gas for transaction")
	cmd.AddCommand(send)

	cmd.SetUsageTemplate(usageTemplate)

	return cmd
}
