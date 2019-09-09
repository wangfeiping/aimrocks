package commands

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/wangfeiping/aimrocks/config"
	"github.com/wangfeiping/aimrocks/log"
)

// NewInitCommand returns init command
func NewInitCommand(run Runner) *cobra.Command {
	cmd := &cobra.Command{
		Use:   CmdInit,
		Short: "Init the node of blockchain",
		Run: func(cmd *cobra.Command, args []string) {
			f := cmd.Flag(FlagConfig)
			if !f.Changed {
				home := viper.GetString(FlagHome)
				configFile := config.
					Check(home, config.DefaultInitConfigFile)
				viper.Set(FlagConfig, configFile)
			}
			if viper.GetBool(FlagCreateInitConfig) {
				createInitConfig()
				return
			}
			loadInitConfig()
			if _, err := run(); err != nil {
				log.Error("chain node init error: ", err)
			}
			return
		},
	}

	cmd.Flags().BoolP(FlagCreateInitConfig, "c",
		false, "Create a new init config file")

	return cmd
}

func createInitConfig() {
	home := viper.GetString(FlagHome)
	configFile := viper.GetString(FlagConfig)
	configFile = config.Check(home, configFile)
	config.EnsureRoot(home)
	config.Create(configFile)
	log.Infof("Config file created: %s", configFile)
}

func loadInitConfig() error {
	home := viper.GetString(FlagHome)
	configFile := viper.GetString(FlagConfig)
	configFile = config.Check(home, configFile)
	config.Load(home, configFile)
	log.Debugf("config file: %s", configFile)
	viper.Set(FlagConfig, configFile)
	return nil
}
