package commands

import (
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func Test_checkConfigFilePath(t *testing.T) {
	viper.Set(FlagHome, "/home/.aimrockscli/")
	viper.Set(FlagConfig, "cli.toml")

	path := checkConfigFilePath()
	assert.Equal(t,
		"/home/.aimrockscli/config/cli.toml",
		path, "wrong path")

	viper.Set(FlagHome, "/home/.aimrockscli/")
	viper.Set(FlagConfig, "./cli.toml")

	path = checkConfigFilePath()
	assert.Equal(t,
		"./cli.toml",
		path, "wrong path")

	viper.Set(FlagHome, "/home/.aimrockscli/")
	viper.Set(FlagConfig, "/root/cli.toml")

	path = checkConfigFilePath()
	assert.Equal(t,
		"/root/cli.toml",
		path, "wrong path")
}
