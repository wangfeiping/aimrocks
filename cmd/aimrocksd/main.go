package main

import (
	"github.com/wangfeiping/aimrocks/commands"
	"github.com/wangfeiping/aimrocks/log"
)

func main() {
	defer log.Flush()

	root := commands.NewRootCommand(versioner)
	root.AddCommand(
		commands.NewStartCommand(nil, true),
		commands.NewAccountCommand(nil, false),
		commands.NewKeyCommand(nil, false),
		commands.NewTxCommand(nil, false),
		commands.NewQueryCommand(nil, false),
		commands.NewVersionCommand(versioner, false))

	if err := root.Execute(); err != nil {
	}
}
