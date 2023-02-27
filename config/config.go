package config

import (
	"strings"
	"time"

	"github.com/spf13/viper"
)

// Config holds the configuration values for the backend.
type Config struct {
	DBmaxopenconns    int
	DBconnmaxlifetime time.Duration
	GitHubToken       string
}

// New loads the config file into the Config struct.
func New() *Config {
	config := viper.New()
	replacer := strings.NewReplacer(".", "_")
	config.SetEnvKeyReplacer(replacer)
	config.AutomaticEnv()

	config.SetDefault("postgres.maxopenconns", 15)
	config.SetDefault("postgres.connmaxlifetime", 15*time.Minute)
	config.SetDefault("github.token", "")

	return &Config{
		DBmaxopenconns:    config.GetInt("postgres.maxopenconns"),
		DBconnmaxlifetime: config.GetDuration("postgres.connmaxlifetime"),
		GitHubToken:       config.GetString("github.token"),
	}
}
