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
	}

	send := &cobra.Command{
		Use:   CmdTxSend,
		Short: "Send transaction(s) to ...",
		Run: func(cmd *cobra.Command, args []string) {
			log.Warn("not implemented yet!")

			// direct to QOS
			// Tx will be sent directly to QOS

			// qstars/bank/processmultitrans.go/MultiSendDirect(...)
			//
			// address:
			//     utility.PubAddrRetrievalFromAmino(...)
			//     types.AccAddressFromBech32(...)
			//     account.AddressStoreKey(...)
			//
			// transactions:
			//     tx.NewTransferMultiple(...)
			//     genStdSendMultiTx(...)
			//
			// submit tx:
			//     cliCtx := *config.GetCLIContext().QOSCliContext
			//     utils.SendTx(...)

			// relay to QOS
			// Tx will be sent to the AimRocksD,
			// and then Cassini will relay the Tx from AimRocksD to QOS

			// MultiSendViaQStars(...)
			//
			// submit tx:
			//     cliCtx := *config.GetCLIContext().QSCCliContext
			//     utils.SendTx(...)

		},
	}

	cmd.AddCommand(send)

	cmd.SetUsageTemplate(usageTemplate)

	return cmd
}
