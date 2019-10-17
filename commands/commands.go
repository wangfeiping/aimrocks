package commands

import (
	"fmt"

	"github.com/spf13/viper"
	"github.com/wangfeiping/aimrocks/config"
)

// nolint
const (
	CmdRootCLI       = "aimrockscli"
	CmdStart         = "start"
	CmdInit          = "init"
	CmdAccount       = "account"
	CmdKey           = "key"
	CmdTx            = "tx"
	CmdTxSend        = "send"
	CmdQuery         = "query"
	CmdVersion       = "version"
	CmdHelp          = "help"
	ShortDescription = "A demo for blockchain"
)

// nolint
const (
	FlagVersion    = CmdVersion
	FlagHome       = "home"
	FlagConfig     = "config"
	FlagLog        = "log"
	FlagNew        = "new"
	FlagFrom       = "from"
	FlagFromAmount = "fromamount"
	FlagTo         = "to"
	FlagToAmount   = "toamount"
	FlagRelay      = "relay"
	FlagTrustNode  = "trust-node"
	FlagMaxGas     = "max-gas"
	FlagCreator    = "creator"
)

var usageTemplate = `Usage:{{if .Runnable}}
  {{.UseLine}}{{end}}{{if .HasAvailableSubCommands}}
  {{.CommandPath}} [command]{{end}}{{if (or (eq .Use "tx") (eq .Use "account"))}}

Query:
  %s %s {{.Use}} [flags]{{end}}{{if gt (len .Aliases) 0}}

Aliases:
  {{.NameAndAliases}}{{end}}{{if .HasExample}}

Examples:
{{.Example}}{{end}}{{if .HasAvailableSubCommands}}

Available Commands:{{range .Commands}}{{if (or .IsAvailableCommand (eq .Name "help"))}}
  {{rpad .Name .NamePadding }} {{.Short}}{{end}}{{end}}{{end}}{{if .HasAvailableLocalFlags}}

Flags:
{{.LocalFlags.FlagUsages | trimTrailingWhitespaces}}{{end}}{{if .HasAvailableInheritedFlags}}

Global Flags:
{{.InheritedFlags.FlagUsages | trimTrailingWhitespaces}}{{end}}{{if .HasHelpSubCommands}}

Additional help topics:{{range .Commands}}{{if .IsAdditionalHelpTopicCommand}}
  {{rpad .CommandPath .CommandPathPadding}} {{.Short}}{{end}}{{end}}{{end}}{{if .HasAvailableSubCommands}}

Use "{{.CommandPath}} [command] --help" for more information about a command.{{end}}
`

var rootCmd string

// func init() {
// 	usageTemplate = fmt.Sprintf(usageTemplate, "aimrocksd", CmdQuery)
// }

// Init the tamplate of usage info
func Init(cmdRoot string) {
	rootCmd = cmdRoot
	usageTemplate = fmt.Sprintf(usageTemplate, cmdRoot, CmdQuery)
}

// GetCmdRoot returns the value of root cmd
func GetCmdRoot() string {
	return rootCmd
}

func checkConfigFilePath() string {
	home := viper.GetString(FlagHome)
	configFile := viper.GetString(FlagConfig)
	configFile = config.Check(home, configFile)
	viper.Set(FlagConfig, configFile)
	return configFile
}
