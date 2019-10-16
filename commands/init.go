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
			f := cmd.Flag(FlagConfig)
			if !f.Changed {
				home := viper.GetString(FlagHome)
				configFile := config.
					Check(home, config.DefaultInitFile)
				viper.Set(FlagConfig, configFile)
			}
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
	configFile = config.Check(home, configFile)
	config.EnsureRoot(home)
	config.Create(configFile)
	log.Infof("init file created: %s", configFile)
}

func loadInitFile() error {
	home := viper.GetString(FlagHome)
	configFile := viper.GetString(FlagConfig)
	configFile = config.Check(home, configFile)
	config.Load(home, configFile)
	log.Debugf("init file loaded: %s", configFile)
	viper.Set(FlagConfig, configFile)
	return nil
}
