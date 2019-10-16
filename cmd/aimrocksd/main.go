package main

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/wangfeiping/aimrocks/commands"
	"github.com/wangfeiping/aimrocks/config"
	"github.com/wangfeiping/aimrocks/log"
)

func main() {
	defer log.Flush()

	// disable sorting
	cobra.EnableCommandSorting = false

	root := commands.NewRootCommand(versioner)
	root.AddCommand(
		commands.NewStartCommand(nil),
		commands.LineBreak,
		commands.NewInitCommand(chainNodeInit),
		commands.NewAccountCommand(nil),
		commands.NewKeyCommand(nil),
		commands.NewTxCommand(txSend),
		commands.NewQueryCommand(nil),
		commands.NewVersionCommand(versioner))

	defaultHome := os.ExpandEnv(config.DefaultHome)
	root.PersistentFlags().String(
		commands.FlagHome,
		defaultHome, "directory for config and data")
	root.PersistentFlags().String(
		commands.FlagConfig,
		config.DefaultConfigFile, "config file path")
	root.PersistentFlags().String(
		commands.FlagLog,
		config.DefaultLogConfigFile, "log config file path")

	if err := root.Execute(); err != nil {
		log.Errorf("Command running error: %v", err)
	}
}
