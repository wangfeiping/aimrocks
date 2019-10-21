package main

import (
	"os"

	"github.com/QOSGroup/qos/app"
	"github.com/spf13/cobra"
	"github.com/wangfeiping/aimrocks/commands"
	"github.com/wangfeiping/aimrocks/config"
	"github.com/wangfeiping/aimrocks/log"
)

func main() {
	defer log.Flush()

	// disable sorting
	cobra.EnableCommandSorting = false
	cdc := app.MakeCodec()
	config.SetDefaultHome(config.DefaultClientHome)
	config.SetDefaultConfigFile(config.DefaultClientFile)
	commands.Init(commands.CmdRootCLI)

	root := commands.NewRootCommand(versioner)
	root.AddCommand(
		commands.NewInitCommand(chainNodeInit),
		commands.NewAccountCommand(nil),
		commands.NewKeysCommand(cdc),
		commands.NewTxCommand(txSend),
		commands.NewQueryCommand(nil),
		commands.NewVersionCommand(versioner))

	defaultHome := os.ExpandEnv(config.GetDefaultHome())
	root.PersistentFlags().String(
		commands.FlagHome,
		defaultHome, "directory for config and data")
	root.PersistentFlags().String(
		commands.FlagConfig,
		config.GetDefaultConfigFile(), "config file path")
	root.PersistentFlags().String(
		commands.FlagLog,
		config.DefaultLogConfigFile, "log config file path")

	if err := root.Execute(); err != nil {
		log.Errorf("Command running error: %v", err)
	}
}
