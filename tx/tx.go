package tx

import (
	qbasetypes "github.com/QOSGroup/qbase/types"
	"github.com/QOSGroup/qstars/types"
	"github.com/QOSGroup/qstars/wire"
	"github.com/wangfeiping/aimrocks/log"
)

// SendTx submit transaction to QOS or QSC
func SendTx(fromAddrs []qbasetypes.Address, fromCoins []types.Coins,
	toAddrs []qbasetypes.Address, toCoins []types.Coins,
	cdc *wire.Codec) {
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

	log.Debug("from addresses: ", len(fromAddrs))
	log.Debug("from coins: ", len(fromCoins))
	log.Debug("to addresses: ", len(toAddrs))
	log.Debug("to coins: ", len(toCoins))
}
