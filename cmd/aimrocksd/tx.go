package main

import (
	"context"
	"fmt"
	"strings"

	qbasetypes "github.com/QOSGroup/qbase/types"
	"github.com/QOSGroup/qstars/star"
	sdk "github.com/QOSGroup/qstars/types"
	"github.com/spf13/viper"
	"github.com/wangfeiping/aimrocks/commands"
	"github.com/wangfeiping/aimrocks/log"
	"github.com/wangfeiping/aimrocks/tx"
)

var txSend = func() (context.CancelFunc, error) {
	from := viper.GetString(commands.FlagFrom)
	famount := viper.GetString(commands.FlagFromAmount)
	to := viper.GetString(commands.FlagTo)
	tamount := viper.GetString(commands.FlagToAmount)

	cdc := star.MakeCodec()

	// parse address
	addrs := []qbasetypes.Address{}
	as := func(str string) error {
		addr, err := sdk.AccAddressFromBech32(str)
		if err != nil {
			return err
		}
		addrs = append(addrs, addr)
		return nil
	}
	if err := parse(as, from); err != nil {
		log.Error("flag parse error: ", err)
		return nil, err
	}
	fromAddrs := addrs

	addrs = addrs[:0]
	if err := parse(as, to); err != nil {
		log.Error("flag parse error: ", err)
		return nil, err
	}

	// parse coins
	coins := []sdk.Coins{}
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
	if err := parse(as, famount); err != nil {
		log.Error("flag parse error: ", err)
		return nil, err
	}
	fromCoins := coins

	coins = coins[:0]
	if err := parse(as, tamount); err != nil {
		log.Error("flag parse error: ", err)
		return nil, err
	}

	tx.SendTx(fromAddrs, fromCoins, addrs, coins, cdc)
	return nil, nil
}

type assembler func(str string) error

func parse(as assembler, str string) (err error) {
	strs := strings.Split(str, ":")
	log.Trace("parsing: ", strs)
	for _, addr := range strs {
		fmt.Println("string: ", addr)
		err = as(addr)
		if err != nil {
			return err
		}
	}
	return nil
}
