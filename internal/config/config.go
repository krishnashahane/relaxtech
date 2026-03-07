package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/viper"
)

// Settings represents the application configuration values loaded from
// file, environment variables, or defaults.
type Settings struct {
	Email        string   `mapstructure:"email"`
	Password     string   `mapstructure:"password"`
	UserID       string   `mapstructure:"user_id"`
	ClientID     string   `mapstructure:"client_id"`
	ClientSecret string   `mapstructure:"client_secret"`
	Timezone     string   `mapstructure:"timezone"`
	Output       string   `mapstructure:"output"`
	Fields       []string `mapstructure:"fields"`
	Verbose      bool     `mapstructure:"verbose"`
}

// ReadSettings reads configuration from the given path (or the default location),
// layers in environment variables prefixed with RELAXTECH_, and returns the
// merged result.
func ReadSettings(cfgFile string, suppressLog bool) (Settings, error) {
	v := viper.New()

	v.SetConfigType("yaml")
	v.SetEnvPrefix("RELAXTECH")
	v.SetEnvKeyReplacer(strings.NewReplacer("-", "_", ".", "_"))
	v.AutomaticEnv()

	if cfgFile != "" {
		v.SetConfigFile(cfgFile)
	} else {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return Settings{}, fmt.Errorf("unable to determine home directory: %w", err)
		}
		v.AddConfigPath(filepath.Join(homeDir, ".config", "relaxtech"))
		v.SetConfigName("config")
	}

	// Apply sensible defaults for optional values.
	v.SetDefault("timezone", "local")
	v.SetDefault("output", "table")

	if err := v.ReadInConfig(); err == nil {
		if !suppressLog {
			fmt.Fprintf(os.Stderr, "Using config file: %s\n", v.ConfigFileUsed())
		}
	}

	var settings Settings
	if err := v.Unmarshal(&settings); err != nil {
		return Settings{}, fmt.Errorf("failed to parse configuration: %w", err)
	}

	return settings, nil
}

// CheckFilePermissions verifies that the configuration file at the given path
// is not readable by group or others. Returns an error with guidance if the
// permissions are too open.
func CheckFilePermissions(filePath string) error {
	if filePath == "" {
		return nil
	}
	stat, err := os.Stat(filePath)
	if err != nil {
		return nil
	}
	perms := stat.Mode().Perm()
	if perms&0o077 != 0 {
		return fmt.Errorf("config file %s has permissions %o; recommended 600", filePath, perms)
	}
	return nil
}
