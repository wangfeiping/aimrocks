package main

import (
	"context"
	"fmt"

	"github.com/wangfeiping/aimrocks/commands"
)

var (
	// Version of cassini
	Version = "0.0.0"

	// GitCommit is the current HEAD set using ldflags.
	GitCommit string

	// GoVersion is version info of golang
	GoVersion string

	// BuidDate is compile date and time
	BuidDate string
)

var versioner = func() (context.CancelFunc, error) {

	s := `AimRocksD - %s
version:	%s
revision:	%s
compile:	%s
go version:	%s
`

	fmt.Printf(s, commands.ShortDescription, Version, GitCommit, BuidDate, GoVersion)

	return nil, nil
}
