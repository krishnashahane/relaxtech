// Package cmd defines the CLI commands for the relaxtech tool.
package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/krishna/relaxtech/internal/client"
	"github.com/krishna/relaxtech/internal/config"
	"github.com/krishna/relaxtech/internal/tokencache"
)

var (
	// primaryCmd is the top-level command that all subcommands attach to.
	primaryCmd = &cobra.Command{
		Use:   "relaxtech",
		Short: "Manage your Eight Sleep Pod directly from the command line",
	}
	// appLogger writes diagnostic messages to stderr.
	appLogger = log.New(os.Stderr)
)

// Execute runs the root command and terminates on failure.
func Execute() {
	if err := primaryCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}

func init() {
	cobra.OnInitialize(loadSettings)

	primaryCmd.PersistentFlags().String("config", "", "settings file (default ~/.config/relaxtech/config.yaml)")
	primaryCmd.PersistentFlags().BoolP("verbose", "v", false, "enable detailed logging")
	primaryCmd.PersistentFlags().String("email", "", "Eight Sleep account email address")
	primaryCmd.PersistentFlags().String("password", "", "Eight Sleep account password")
	primaryCmd.PersistentFlags().String("client-id", "", "Eight Sleep client ID (optional; falls back to public app client)")
	primaryCmd.PersistentFlags().String("client-secret", "", "Eight Sleep client secret (optional; falls back to public app client)")
	primaryCmd.PersistentFlags().String("user-id", "", "Eight Sleep user identifier")
	primaryCmd.PersistentFlags().String("timezone", "local", "IANA timezone name (e.g., America/New_York) or 'local'")
	primaryCmd.PersistentFlags().String("output", "table", "display mode: table|json|csv")
	primaryCmd.PersistentFlags().StringSlice("fields", []string{}, "restrict output to these columns")
	primaryCmd.PersistentFlags().Bool("quiet", false, "hide settings load messages")

	viper.BindPFlag("config", primaryCmd.PersistentFlags().Lookup("config"))
	viper.BindPFlag("verbose", primaryCmd.PersistentFlags().Lookup("verbose"))
	viper.BindPFlag("email", primaryCmd.PersistentFlags().Lookup("email"))
	viper.BindPFlag("password", primaryCmd.PersistentFlags().Lookup("password"))
	viper.BindPFlag("client_id", primaryCmd.PersistentFlags().Lookup("client-id"))
	viper.BindPFlag("client_secret", primaryCmd.PersistentFlags().Lookup("client-secret"))
	viper.BindPFlag("user_id", primaryCmd.PersistentFlags().Lookup("user-id"))
	viper.BindPFlag("timezone", primaryCmd.PersistentFlags().Lookup("timezone"))
	viper.BindPFlag("output", primaryCmd.PersistentFlags().Lookup("output"))
	viper.BindPFlag("fields", primaryCmd.PersistentFlags().Lookup("fields"))
	viper.BindPFlag("config-quiet", primaryCmd.PersistentFlags().Lookup("quiet"))

	primaryCmd.AddCommand(activateCmd)
	primaryCmd.AddCommand(deactivateCmd)
	primaryCmd.AddCommand(temperatureCmd)
	primaryCmd.AddCommand(deviceStatusCmd)
	primaryCmd.AddCommand(buildVersionCmd)
	primaryCmd.AddCommand(identityCmd)
	primaryCmd.AddCommand(signoutCmd)
	primaryCmd.AddCommand(sleepCmd)
	primaryCmd.AddCommand(alarmCmd)
	primaryCmd.AddCommand(scheduleCmd)
	primaryCmd.AddCommand(audioCmd)
	primaryCmd.AddCommand(tracksCmd)
	primaryCmd.AddCommand(deviceCmd)
	primaryCmd.AddCommand(baseCmd)
	primaryCmd.AddCommand(metricsCmd)
	primaryCmd.AddCommand(autopilotCmd)
	primaryCmd.AddCommand(householdCmd)
	primaryCmd.AddCommand(travelCmd)
	primaryCmd.AddCommand(presenceCmd)
	primaryCmd.AddCommand(featsCmd)
	primaryCmd.AddCommand(tempModeCmd)
	primaryCmd.AddCommand(daemonCmd)
}

// loadSettings reads the configuration file and applies values to viper.
func loadSettings() {
	settings, err := config.ReadSettings(viper.GetString("config"), viper.GetBool("config-quiet"))
	if err != nil {
		log.Fatalf("settings: %v", err)
	}

	// Allow environment variables to override configuration values.
	viper.SetEnvPrefix("RELAXTECH")
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_", ".", "_"))
	viper.AutomaticEnv()

	// Populate viper defaults from the loaded settings.
	viper.SetDefault("email", settings.Email)
	viper.SetDefault("password", settings.Password)
	viper.SetDefault("user_id", settings.UserID)
	viper.SetDefault("client_id", settings.ClientID)
	viper.SetDefault("client_secret", settings.ClientSecret)
	viper.SetDefault("timezone", settings.Timezone)
	viper.SetDefault("output", settings.Output)
	viper.SetDefault("fields", settings.Fields)
	viper.SetDefault("verbose", settings.Verbose)

	if err := config.CheckFilePermissions(viper.ConfigFileUsed()); err != nil {
		appLogger.Warn(err.Error())
	}

	if viper.GetBool("verbose") {
		log.SetLevel(log.DebugLevel)
	}
}

// ensureCredentials checks that authentication details are available,
// either from a stored token or from explicit credentials.
func ensureCredentials() error {
	// First try the token cache to avoid requiring raw credentials.
	cl := client.New(
		viper.GetString("email"),
		viper.GetString("password"),
		viper.GetString("user_id"),
		viper.GetString("client_id"),
		viper.GetString("client_secret"),
	)
	if stored, err := tokencache.Retrieve(cl.AuthContext(), viper.GetString("user_id")); err == nil {
		if stored.UserID != "" {
			viper.Set("user_id", stored.UserID)
		}
		return nil
	}

	// Fall back to checking that credentials were provided.
	absent := []string{}
	if viper.GetString("email") == "" {
		absent = append(absent, "email")
	}
	if viper.GetString("password") == "" {
		absent = append(absent, "password")
	}
	if len(absent) > 0 {
		return fmt.Errorf("missing required auth fields: %s", strings.Join(absent, ", "))
	}
	return nil
}
