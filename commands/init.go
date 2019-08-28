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
			if viper.GetBool(FlagCreateConfig) {
				createConfig()
				return
			}

			if _, err := run(); err != nil {
				log.Error("chain node init error: ", err)
			}
			return
		},
	}

	cmd.Flags().BoolP(FlagCreateConfig, "c",
		false, "Create a new config file")

	return cmd
}

func createConfig() {
	home := viper.GetString(FlagHome)
	configFile := viper.GetString(FlagConfig)
	configFile = config.Check(home, configFile)
	config.EnsureRoot(home)
	config.Create(configFile)
	log.Infof("Config file created: %s", configFile)
}
