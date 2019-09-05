package commands

import (
	"context"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/cihub/seelog"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/wangfeiping/aimrocks/config"
	"github.com/wangfeiping/aimrocks/log"
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

			if strings.EqualFold(cmd.Use, CmdRoot) ||
				strings.EqualFold(cmd.Use, CmdVersion) ||
				(strings.EqualFold(cmd.Use, CmdInit) &&
					viper.GetBool(FlagCreateInitConfig)) {
				// doesn't need init config & log
				return nil
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
	home := viper.GetString(FlagHome)
	configFile := viper.GetString(FlagConfig)
	configFile = config.Check(home, configFile)
	config.Load(home, configFile)
	log.Debugf("config file: %s", configFile)
	viper.Set(FlagConfig, configFile)
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
