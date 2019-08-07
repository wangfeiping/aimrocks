package main

import (
	"github.com/spf13/cobra"
	"github.com/wangfeiping/aimrocks/commands"
	"github.com/wangfeiping/aimrocks/log"
)

func main() {
	defer log.Flush()

	// disable sorting
	cobra.EnableCommandSorting = false

	root := commands.NewRootCommand(versioner)
	root.AddCommand(
		commands.NewStartCommand(nil, true),
		commands.NewInitCommand(nil, false),
		commands.NewAccountCommand(nil, false),
		commands.NewKeyCommand(nil, false),
		commands.NewTxCommand(txSend),
		commands.NewQueryCommand(nil, false),
		commands.NewVersionCommand(versioner))

	if err := root.Execute(); err != nil {
	}
}
