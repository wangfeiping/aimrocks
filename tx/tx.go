package tx

import (
	qbasetypes "github.com/QOSGroup/qbase/types"
	"github.com/QOSGroup/qstars/types"
	"github.com/QOSGroup/qstars/wire"
	"github.com/QOSGroup/qstars/x/bank"
	"github.com/spf13/viper"
	"github.com/wangfeiping/aimrocks/commands"
	"github.com/wangfeiping/aimrocks/log"
)

// SendTx submit transaction to QOS or QSC
func SendTx(fromAddrs []qbasetypes.Address, fromCoins []types.Coins,
	toAddrs []qbasetypes.Address, toCoins []types.Coins,
	cdc *wire.Codec) (*bank.SendResult, error) {
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
	//     genStdSendMultiTx(...) !!! max-gas
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

	from := []string{"mWMFi5yQeEghGgL+dL4zop4US45lvZ1Ag4Ai5J8041Fj3RqLnUtGtc5AVQ0gxFn3enFkWcCTmCw3IH7rFezjkg=="}
	var result *bank.SendResult
	var err error
	// !!! must change max-gas code !!!
	if viper.GetBool(commands.FlagRelay) {
		log.Debug("tx send relay")
		result, err = bank.MultiSendViaQStars(
			cdc, from, toAddrs, fromCoins, toCoins)
		if err != nil {
			return nil, err
		}
	} else {
		log.Debug("tx send direct")
		result, err = bank.MultiSendDirect(
			cdc, from, toAddrs, fromCoins, toCoins)
		if err != nil {
			return nil, err
		}
	}
	return result, err
}
