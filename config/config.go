package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/viper"
	"github.com/wangfeiping/aimrocks/log"
)

// nolint
const (
	DefaultHome          = "$HOME/.aimrocks"
	DefaultConfigFile    = "config.toml"
	DefaultLogConfigFile = "log.conf"
)

// Config wraps configure data
type Config struct {
	LogConfigFile       string `mapstructure:"log"`
	QOSChainID          string `mapstructure:"qos_chain_id"`
	QSCChainID          string `mapstructure:"qsc_chain_id"`
	QOSNodeURI          string `mapstructure:"qos_node_uri"`
	QSTARSNodeURI       string `mapstructure:"qstars_node_uri"`
	DirectTOQOS         bool   `mapstructure:"direct_to_qos"`
	WaitingForQosResult string `mapstructure:"waiting_for_qos_result"`
	Kepler              string `mapstructure:"kepler"`

	Community   string `mapstructure:"community"`
	Authormock  string `mapstructure:"authormock"`
	Adbuyermock string `mapstructure:"adbuyermock"`
	Banker      string `mapstructure:"banker"`
	Dappowner   string `mapstructure:"dappowner"`
}

var conf = DefaultConfig()

// GetConfig returns the config instance
func GetConfig() *Config {
	return conf
}

// DefaultConfig creates a default config
func DefaultConfig() *Config {
	c := &Config{
		LogConfigFile: DefaultLogConfigFile,
		QSCChainID:    "qstars-test",
		QOSChainID:    "qos-test",
		QOSNodeURI:    "localhost:26657",
		QSTARSNodeURI: "localhost:26657",
		Kepler:        "http://127.0.0.1:8080/kepler/"}
	return c
}

// Check returns actual config file path
func Check(home, configFile string) string {
	configDir := filepath.Join(home, "config")

	dir, file := filepath.Split(configFile)
	if dir == "" {
		dir = configDir
	}
	if file == "" {
		file = DefaultConfigFile
	}
	file = filepath.Join(dir, file)
	if strings.EqualFold(file, DefaultConfigFile) {
		return configFile
	}
	return file
}

// Create writes a new config file
func Create(configFilePath string) {
	// Create only if the file doesn't exist
	if _, err := os.Stat(configFilePath); os.IsNotExist(err) { //
		// the following parse config is needed to create directories
		conf = DefaultConfig()
		WriteConfigFile(configFilePath, conf)
		// Fall through, just so that its parsed into memory.
	} else if err != nil {
		panic(err)
	} else {
		panic(fmt.Errorf("Config file is exist: %s", configFilePath))
	}
	return
}

// Load loads config data from file
func Load(home, configFile string) {
	log.Debug("home: ", home)
	log.Debug("config: ", configFile)

	viper.SetConfigFile(configFile)

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err != nil {
		log.Errorf("Load config(%s) failed: %v",
			configFile, err)
		panic(err)
	}
}
