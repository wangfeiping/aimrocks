package main

import (
	"context"

	"github.com/spf13/viper"
	"github.com/wangfeiping/aimrocks/commands"
	"github.com/wangfeiping/aimrocks/log"
)

var chainNodeInit = func() (context.CancelFunc, error) {
	home := viper.GetString(commands.FlagHome)
	config := viper.GetString(commands.FlagConfig)
	log.Infof("chain node init... %s", home)
	log.Infof("chain node init... %s", config)
	return nil, nil
}
