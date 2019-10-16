package commands

import "fmt"

// nolint
const (
	CmdRoot          = "aimrocksd"
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
  {{.CommandPath}} [command]{{end}}{{if (or (eq .Use "tx") (eq .Use "account") (eq .Use "key"))}}

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

func init() {
	usageTemplate = fmt.Sprintf(usageTemplate, CmdRoot, CmdQuery)
}
