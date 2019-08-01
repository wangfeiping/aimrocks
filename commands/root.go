package commands

import (
	"context"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/cihub/seelog"
	"github.com/spf13/cobra"
	"github.com/wangfeiping/aimrocks/config"
	"github.com/wangfeiping/aimrocks/log"
)

const (
	// CmdRoot string of root command
	CmdRoot = "aimrocksd"

	// CmdVersion string of version command
	CmdVersion = "version"

	// CmdHelp string of help command
	CmdHelp = "help"

	// CmdStart string of start command
	CmdStart = "start"

	// CmdAccount string of account command
	CmdAccount = "account"

	// CmdTx string of tx command
	CmdTx = "tx"

	// CmdQuery string of query command
	CmdQuery = "query"

	// ShortDescription string of short description
	ShortDescription = "A demo for blockchain"
)

// Runner is command call function
type Runner func() (context.CancelFunc, error)

// NewRootCommand returns root command
func NewRootCommand(versioner Runner) *cobra.Command {
	root := &cobra.Command{
		Use:   CmdRoot,
		Short: ShortDescription,
		Run: func(cmd *cobra.Command, args []string) {
			if config.GetConfig().ShowVersion {
				versioner()
				return
			}
			cmd.Help()
		},
		PersistentPreRunE: func(cmd *cobra.Command, args []string) (err error) {
			// if strings.EqualFold(cmd.Use, CmdRoot) ||
			// 	strings.EqualFold(cmd.Use, CmdVersion) ||
			// 	strings.HasPrefix(cmd.Use, CmdHelp) {
			if !strings.EqualFold(cmd.Use, CmdStart) {
				// doesn't need init log and config
				return nil
			}
			// init logger
			var logger seelog.LoggerInterface
			logger, err = log.LoadLogger(config.GetConfig().LogConfigFile)
			if err != nil {
				log.Warn("Used the default logger because error: ", err)
			} else {
				log.Replace(logger)
			}
			// err = initConfig()
			// if err != nil {
			// 	return err
			// }
			return
		},
	}
	// root.Version = "asd"
	root.Flags().BoolVarP(&config.GetConfig().ShowVersion,
		"version", "v", false, "Show version info")
	return root
}

func commandRunner(run Runner, isKeepRunning bool) error {
	cancel, err := run()
	if err != nil {
		log.Error("Run command error: ", err.Error())
		return err
	}
	if isKeepRunning {
		keepRunning(func(sig os.Signal) {
			defer log.Flush()
			if cancel != nil {
				cancel()
			}
			log.Debug("Stopped by signal: ", sig)
		})
	}
	return nil
}

func keepRunning(callback func(sig os.Signal)) {
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt, syscall.SIGTERM)

	select {
	case s, ok := <-signals:
		log.Infof("System signal [%v] %t, trying to run callback...", s, ok)
		if !ok {
			break
		}
		if callback != nil {
			callback(s)
		}
		log.Flush()
		os.Exit(1)
	}
}
