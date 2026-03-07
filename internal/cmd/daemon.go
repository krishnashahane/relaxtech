package cmd

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"

	"github.com/krishna/relaxtech/internal/client"
	"github.com/krishna/relaxtech/internal/daemon"
)

// daemonCmd launches the relaxtech schedule daemon, which reads timed tasks
// from the configuration file and executes them on their schedule.
var daemonCmd = &cobra.Command{
	Use:   "daemon",
	Short: "Run schedule daemon from config file",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := ensureCredentials(); err != nil {
			return err
		}
		configBytes, err := loadSettingsSchedule()
		if err != nil {
			return err
		}
		tasks, err := extractTimedTasks(configBytes)
		if err != nil {
			return err
		}
		tzLabel := viper.GetString("timezone")
		region := time.Local
		if tzLabel != "local" {
			region, err = time.LoadLocation(tzLabel)
			if err != nil {
				return fmt.Errorf("load timezone: %w", err)
			}
		}
		sched := daemon.Scheduler{
			Tasks:     tasks,
			DeviceAPI: client.New(viper.GetString("email"), viper.GetString("password"), viper.GetString("user_id"), viper.GetString("client_id"), viper.GetString("client_secret")),
			Location:  region,
			DryRun:   viper.GetBool("dry-run"),
			Sync:     viper.GetBool("sync-state"),
			PIDFilePath: resolvePIDPath(viper.GetString("pid-file")),
		}
		ctx := context.Background()
		fmt.Printf("daemon started with %d items\n", len(tasks))
		return sched.Start(ctx)
	},
}

func init() {
	daemonCmd.Flags().Bool("dry-run", false, "log actions without executing")
	daemonCmd.Flags().Bool("sync-state", false, "(reserved) sync device state")
	daemonCmd.Flags().String("pid-file", "", "pid file path (default ~/.config/relaxtech/daemon.pid)")
	viper.BindPFlag("dry-run", daemonCmd.Flags().Lookup("dry-run"))
	viper.BindPFlag("sync-state", daemonCmd.Flags().Lookup("sync-state"))
	viper.BindPFlag("pid-file", daemonCmd.Flags().Lookup("pid-file"))
}

// loadSettingsSchedule reads the raw bytes from the active configuration file
// so that the schedule section can be parsed independently.
func loadSettingsSchedule() ([]byte, error) {
	settingsPath := viper.ConfigFileUsed()
	if settingsPath == "" {
		return nil, fmt.Errorf("no config file loaded; specify --config")
	}
	return os.ReadFile(settingsPath)
}

// extractTimedTasks unmarshals the schedule section of the configuration
// into a slice of TimedTask entries that the Scheduler can process.
func extractTimedTasks(rawData []byte) ([]daemon.TimedTask, error) {
	var wrapper struct {
		Schedule []daemon.TimedTask `yaml:"schedule"`
	}
	if err := yaml.Unmarshal(rawData, &wrapper); err != nil {
		return nil, err
	}
	if len(wrapper.Schedule) == 0 {
		return nil, fmt.Errorf("no schedule entries found")
	}
	return wrapper.Schedule, nil
}

// resolvePIDPath returns the user-specified PID file path, or falls back
// to the default location under the relaxtech config directory.
func resolvePIDPath(userValue string) string {
	if userValue != "" {
		return userValue
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return ""
	}
	return filepath.Join(home, ".config", "relaxtech", "daemon.pid")
}
