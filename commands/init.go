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
		Short: "init the node of blockchain",
		Run: func(cmd *cobra.Command, args []string) {
			if viper.GetBool(FlagNew) {
				createInitFile()
				return
			}
			loadInitFile()
			if _, err := run(); err != nil {
				log.Error("chain node init error: ", err)
			}
			return
		},
	}

	cmd.Flags().BoolP(FlagNew, "n",
		false, "generate a new init file")
	cmd.Flags().BoolP(FlagCreator, "c",
		false, "creator's account for the qcp chain")
	return cmd
}

func createInitFile() {
	home := viper.GetString(FlagHome)
	configFile := viper.GetString(FlagConfig)
	config.EnsureRoot(home)
	config.Create(configFile)
	log.Infof("init file created: %s", configFile)
}

func loadInitFile() error {
	configFile := viper.GetString(FlagConfig)
	config.Load(configFile)
	viper.Set(FlagConfig, configFile)
	log.Debugf("init file loaded: %s", configFile)
	return nil
}
