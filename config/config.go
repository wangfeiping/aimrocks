package config

// Config wraps configure data
type Config struct {

	// ConfigFile path of config file
	ConfigFile string

	// LogConfigFile path of log config file
	LogConfigFile string
}

var conf = &Config{}

// GetConfig returns the config instance
func GetConfig() *Config {
	return conf
}
