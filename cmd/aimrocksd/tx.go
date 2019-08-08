package main

import (
	"context"
	"strings"

	qbasetypes "github.com/QOSGroup/qbase/types"
	"github.com/QOSGroup/qstars/config"
	"github.com/QOSGroup/qstars/star"
	sdk "github.com/QOSGroup/qstars/types"
	"github.com/QOSGroup/qstars/wire"
	"github.com/QOSGroup/qstars/x/bank"
	"github.com/spf13/viper"
	"github.com/wangfeiping/aimrocks/commands"
	"github.com/wangfeiping/aimrocks/log"
)

var txSend = func() (context.CancelFunc, error) {
	fromAddrs, fromCoins, addrs, coins, err :=
		parseTxSendFlags()
	if err != nil {
		log.Error("flags parse error: ", err)
		return nil, err
	}
	log.Debug("from addrs: ", len(fromAddrs))

	cdc := star.MakeCodec()

	cfg := &config.CLIConfig{
		QSCChainID:    "dawns-3001",
		QOSChainID:    "capricorn-3000",
		QOSNodeURI:    "localhost:26657",
		QSTARSNodeURI: "localhost:26658"}
	config.CreateCLIContextTwo(cdc, cfg)

	from := []string{"oSXr2kEsWgw8L9ydeLOLM9Q8A6g+HhMW5sAdrKkycdDoiLLNrfCkR5P+7gNiYDSuyW390yhbnCv4+PXhhf/O0w=="}
	var result *bank.SendResult
	// !!! must change max-gas code !!!
	result, err = bank.MultiSendDirect(cdc,
		from, addrs, fromCoins, coins)
	if err != nil {
		log.Error("tx send error: ", err)
		return nil, err
	}

	// tx.SendTx(fromAddrs, fromCoins, addrs, coins, cdc)

	output, err := wire.MarshalJSONIndent(cdc, result)
	if err != nil {
		log.Error("tx send result parse error: ", err)
		return nil, err
	}
	log.Debugf("tx send result: ", string(output))
	return nil, nil
}

func parseTxSendFlags() (fromAddrs []qbasetypes.Address, fromCoins []sdk.Coins,
	addrs []qbasetypes.Address, coins []sdk.Coins, err error) {
	from := viper.GetString(commands.FlagFrom)
	famount := viper.GetString(commands.FlagFromAmount)
	to := viper.GetString(commands.FlagTo)
	tamount := viper.GetString(commands.FlagToAmount)

	// parse address
	as := func(str string) error {
		addr, err := sdk.AccAddressFromBech32(str)
		if err != nil {
			return err
		}
		addrs = append(addrs, addr)
		return nil
	}
	if err = parse(as, from); err != nil {
		return
	}
	fromAddrs = addrs

	addrs = addrs[:0]
	if err = parse(as, to); err != nil {
		return
	}

	// parse coins
	as = func(str string) error {
		coin, err := sdk.ParseCoins(str)
		if err != nil {
			return err
		}
		if coin == nil {
			return nil
		}
		coins = append(coins, coin)
		return nil
	}
	if err = parse(as, famount); err != nil {
		return
	}
	fromCoins = coins

	coins = coins[:0]
	err = parse(as, tamount)
	return
}

type assembler func(str string) error

func parse(as assembler, str string) (err error) {
	var strs []string
	if strings.Index(str, ":") >= 0 {
		strs = strings.Split(str, ":")
	} else {
		strs = strings.Split(str, ",")
	}
	log.Debug("parsing: ", len(strs))
	for _, addr := range strs {
		log.Debug("string: ", addr)
		err = as(addr)
		if err != nil {
			return err
		}
	}
	return nil
}
