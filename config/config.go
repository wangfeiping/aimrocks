package config

// Config wraps configure data
type Config struct {

	// ConfigFile path of config file
	ConfigFile string

	// LogConfigFile path of log config file
	LogConfigFile string

	// ShowVersion flag that show version info
	ShowVersion bool
}

var conf = &Config{}

// GetConfig returns the config instance
func GetConfig() *Config {
	return conf
}
