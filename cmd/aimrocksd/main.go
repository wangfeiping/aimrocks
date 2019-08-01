package main

import (
	"github.com/wangfeiping/aimrocks/commands"
	"github.com/wangfeiping/aimrocks/log"
)

func main() {
	defer log.Flush()

	// cli := commands.NewClientCommand()
	// cli.AddCommand(
	// 	commands.NewQueryCommand(mocker, true),
	// 	commands.NewAccountCommand(resetHandler, false),
	// 	commands.NewTxCommand(txHandler, false))

	root := commands.NewRootCommand(versioner)
	root.AddCommand(
		commands.NewStartCommand(nil, true),
		commands.NewAccountCommand(nil, true),
		commands.NewTxCommand(nil, true),
		commands.NewQueryCommand(nil, true),
		commands.NewVersionCommand(versioner, false))

	if err := root.Execute(); err != nil {
		log.Error("Exit with error: ", err)
	}
}
