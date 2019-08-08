package commands

import (
	"context"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"

	qstarconfig "github.com/QOSGroup/qstars/config"
	"github.com/QOSGroup/qstars/star"
	"github.com/cihub/seelog"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/wangfeiping/aimrocks/config"
	"github.com/wangfeiping/aimrocks/log"
)

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
	FlagVersion    = "version"
	FlagHome       = "home"
	FlagFrom       = "from"
	FlagFromAmount = "fromamount"
	FlagTo         = "to"
	FlagToAmount   = "toamount"
	FlagRelay      = "relay"
	FlagTrustNode  = "trust-node"
	FlagMaxGas     = "max-gas"
)

// Runner is command call function
type Runner func() (context.CancelFunc, error)

// NewRootCommand returns root command
func NewRootCommand(versioner Runner) *cobra.Command {
	root := &cobra.Command{
		Use:   CmdRoot,
		Short: ShortDescription,
		Run: func(cmd *cobra.Command, args []string) {
			if viper.GetBool(FlagVersion) {
				versioner()
				return
			}
			cmd.Help()
		},
		PersistentPreRunE: func(cmd *cobra.Command, args []string) (err error) {
			if err := viper.BindPFlags(cmd.Flags()); err != nil {
				log.Error("bind flags error: ", err)
				return err
			}

			initConfig()

			if !strings.EqualFold(cmd.Use, CmdStart) {
				// doesn't need init log
				return nil
			}

			initLogger()

			return
		},
	}

	root.Flags().BoolP("version", "v", false, "Show version info")

	return root
}

func initConfig() error {
	// cfg := &qstarconfig.CLIConfig{
	// 	QSCChainID:    "dawns-3001",
	// 	QOSChainID:    "capricorn-3000",
	// 	QOSNodeURI:    "localhost:26657",
	// 	QSTARSNodeURI: "localhost:26658"}
	// config.CreateCLIContextTwo(cdc, cfg)

	log.Debug("home: ", viper.GetString("home"))

	homeDir := viper.GetString(FlagHome)
	viper.Set(FlagHome, homeDir)
	// Sets name for the config file.
	// Does not include extension.
	viper.SetConfigName("config")
	// Adds a path for Viper to search for the config file in.
	viper.AddConfigPath(filepath.Join(homeDir, "config"))
	// Can be called multiple times to define multiple search paths.
	viper.AddConfigPath(homeDir)

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		// stderr, so if we redirect output to json file, this doesn't appear
		// fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	} else if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
		// ignore not found error, return other errors
		return err
	}

	cfg, err := qstarconfig.InterceptLoadConfig()
	if err != nil {
		log.Error("init config error: ", err)
		return err
	}

	log.Debug("QOSNodeURI: ", cfg.QOSNodeURI)
	log.Debug("QSTARSNodeURI: ", cfg.QSTARSNodeURI)
	cdc := star.MakeCodec()
	qstarconfig.CreateCLIContextTwo(cdc, cfg)
	return nil
}

func initLogger() {
	var logger seelog.LoggerInterface
	logger, err := log.LoadLogger(config.GetConfig().LogConfigFile)
	if err != nil {
		log.Warn("Used the default logger because error: ", err)
	} else {
		log.Replace(logger)
	}
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
