package cmd

import (
	"context"
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/krishna/relaxtech/internal/client"
	"github.com/krishna/relaxtech/internal/output"
)

// alarmCmd groups all alarm-related subcommands under the "alarm" namespace.
var alarmCmd = &cobra.Command{
	Use:   "alarm",
	Short: "Manage alarms",
}

// alarmListCmd retrieves and displays all configured alarms.
var alarmListCmd = &cobra.Command{
	Use:   "list",
	Short: "List alarms",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := ensureCredentials(); err != nil {
			return err
		}
		cl := client.New(viper.GetString("email"), viper.GetString("password"), viper.GetString("user_id"), viper.GetString("client_id"), viper.GetString("client_secret"))
		alarmEntries, err := cl.ListAlarms(context.Background())
		if err != nil {
			return err
		}
		records := make([]map[string]any, 0, len(alarmEntries))
		for _, entry := range alarmEntries {
			records = append(records, map[string]any{
				"id":        entry.ID,
				"time":      entry.Time,
				"enabled":   entry.Enabled,
				"days":      entry.DaysOfWeek,
				"vibration": entry.Vibration,
				"sound":     entry.Sound,
			})
		}
		records = output.SelectColumns(records, viper.GetStringSlice("fields"))
		return output.Render(output.DisplayMode(viper.GetString("output")), []string{"id", "time", "enabled", "days", "vibration", "sound"}, records)
	},
}

// alarmCreateCmd adds a new alarm with the specified parameters.
var alarmCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create an alarm",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := ensureCredentials(); err != nil {
			return err
		}
		alarmTime := viper.GetString("time")
		if alarmTime == "" {
			return fmt.Errorf("--time required")
		}
		weekdays := viper.GetIntSlice("days")
		if len(weekdays) == 0 {
			return fmt.Errorf("--days required (comma separated 0=Sun..6=Sat)")
		}
		soundID := viper.GetString("sound")
		var soundRef *string
		if soundID != "" {
			soundRef = &soundID
		}
		cl := client.New(viper.GetString("email"), viper.GetString("password"), viper.GetString("user_id"), viper.GetString("client_id"), viper.GetString("client_secret"))
		newAlarm := client.Alarm{
			Enabled:    !viper.GetBool("disabled"),
			Time:       alarmTime,
			DaysOfWeek: weekdays,
			Vibration:  !viper.GetBool("no-vibration"),
			Sound:      soundRef,
		}
		created, err := cl.CreateAlarm(context.Background(), newAlarm)
		if err != nil {
			return err
		}
		fmt.Printf("created alarm %s\n", created.ID)
		return nil
	},
}

// alarmUpdateCmd modifies an existing alarm identified by its ID.
var alarmUpdateCmd = &cobra.Command{
	Use:   "update <id>",
	Short: "Update an alarm",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := ensureCredentials(); err != nil {
			return err
		}
		changes := map[string]any{}
		if val := viper.GetString("time"); val != "" {
			changes["time"] = val
		}
		if weekdays := viper.GetIntSlice("days"); len(weekdays) > 0 {
			changes["daysOfWeek"] = weekdays
		}
		if cmd.Flags().Changed("enabled") {
			changes["enabled"] = viper.GetBool("enabled")
		}
		if cmd.Flags().Changed("no-vibration") {
			changes["vibration"] = !viper.GetBool("no-vibration")
		}
		if soundID := viper.GetString("sound"); soundID != "" {
			changes["sound"] = soundID
		}
		if len(changes) == 0 {
			return fmt.Errorf("no fields to update")
		}
		cl := client.New(viper.GetString("email"), viper.GetString("password"), viper.GetString("user_id"), viper.GetString("client_id"), viper.GetString("client_secret"))
		if _, err := cl.UpdateAlarm(context.Background(), args[0], changes); err != nil {
			return err
		}
		fmt.Println("updated")
		return nil
	},
}

// alarmDeleteCmd removes an alarm by its ID.
var alarmDeleteCmd = &cobra.Command{
	Use:   "delete <id>",
	Short: "Delete an alarm",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := ensureCredentials(); err != nil {
			return err
		}
		cl := client.New(viper.GetString("email"), viper.GetString("password"), viper.GetString("user_id"), viper.GetString("client_id"), viper.GetString("client_secret"))
		if err := cl.DeleteAlarm(context.Background(), args[0]); err != nil {
			return err
		}
		fmt.Println("deleted")
		return nil
	},
}

func init() {
	alarmCreateCmd.Flags().String("time", "", "HH:MM time")
	alarmCreateCmd.Flags().IntSlice("days", nil, "Comma-separated days 0=Sun..6=Sat")
	alarmCreateCmd.Flags().Bool("disabled", false, "Create disabled")
	alarmCreateCmd.Flags().Bool("no-vibration", false, "Disable vibration")
	alarmCreateCmd.Flags().String("sound", "", "Sound id")
	viper.BindPFlag("time", alarmCreateCmd.Flags().Lookup("time"))
	viper.BindPFlag("days", alarmCreateCmd.Flags().Lookup("days"))
	viper.BindPFlag("disabled", alarmCreateCmd.Flags().Lookup("disabled"))
	viper.BindPFlag("no-vibration", alarmCreateCmd.Flags().Lookup("no-vibration"))
	viper.BindPFlag("sound", alarmCreateCmd.Flags().Lookup("sound"))

	alarmUpdateCmd.Flags().String("time", "", "HH:MM time")
	alarmUpdateCmd.Flags().IntSlice("days", nil, "Comma-separated days 0=Sun..6=Sat")
	alarmUpdateCmd.Flags().Bool("enabled", true, "Set enabled true/false")
	alarmUpdateCmd.Flags().Bool("no-vibration", false, "Disable vibration")
	alarmUpdateCmd.Flags().String("sound", "", "Sound id")
	viper.BindPFlag("time", alarmUpdateCmd.Flags().Lookup("time"))
	viper.BindPFlag("days", alarmUpdateCmd.Flags().Lookup("days"))
	viper.BindPFlag("enabled", alarmUpdateCmd.Flags().Lookup("enabled"))
	viper.BindPFlag("no-vibration", alarmUpdateCmd.Flags().Lookup("no-vibration"))
	viper.BindPFlag("sound", alarmUpdateCmd.Flags().Lookup("sound"))

	// Register all alarm subcommands with the parent alarm command.
	alarmCmd.AddCommand(alarmListCmd, alarmCreateCmd, alarmUpdateCmd, alarmDeleteCmd, alarmSnoozeCmd, alarmDismissCmd, alarmDismissAllCmd, alarmVibeCmd)
}

// alarmSnoozeCmd delays a ringing alarm for a brief period.
var alarmSnoozeCmd = &cobra.Command{Use: "snooze <id>", Args: cobra.ExactArgs(1), RunE: func(cmd *cobra.Command, args []string) error {
	if err := ensureCredentials(); err != nil {
		return err
	}
	cl := client.New(viper.GetString("email"), viper.GetString("password"), viper.GetString("user_id"), viper.GetString("client_id"), viper.GetString("client_secret"))
	return cl.Alarms().Snooze(context.Background(), args[0])
}}

// alarmDismissCmd silences a single active alarm.
var alarmDismissCmd = &cobra.Command{Use: "dismiss <id>", Args: cobra.ExactArgs(1), RunE: func(cmd *cobra.Command, args []string) error {
	if err := ensureCredentials(); err != nil {
		return err
	}
	cl := client.New(viper.GetString("email"), viper.GetString("password"), viper.GetString("user_id"), viper.GetString("client_id"), viper.GetString("client_secret"))
	return cl.Alarms().Dismiss(context.Background(), args[0])
}}

// alarmDismissAllCmd silences every active alarm at once.
var alarmDismissAllCmd = &cobra.Command{Use: "dismiss-all", RunE: func(cmd *cobra.Command, args []string) error {
	if err := ensureCredentials(); err != nil {
		return err
	}
	cl := client.New(viper.GetString("email"), viper.GetString("password"), viper.GetString("user_id"), viper.GetString("client_id"), viper.GetString("client_secret"))
	return cl.Alarms().DismissAll(context.Background())
}}

// alarmVibeCmd triggers a short vibration test through the device.
var alarmVibeCmd = &cobra.Command{Use: "vibration-test", RunE: func(cmd *cobra.Command, args []string) error {
	if err := ensureCredentials(); err != nil {
		return err
	}
	cl := client.New(viper.GetString("email"), viper.GetString("password"), viper.GetString("user_id"), viper.GetString("client_id"), viper.GetString("client_secret"))
	return cl.Alarms().VibrationTest(context.Background())
}}

// splitDayValues splits a comma-separated string of day numbers into integers.
// Kept for potential future use by relaxtech commands.
func splitDayValues(s string) ([]int, error) {
	segments := strings.Split(s, ",")
	result := make([]int, 0, len(segments))
	for _, seg := range segments {
		seg = strings.TrimSpace(seg)
		if seg == "" {
			continue
		}
		var dayNum int
		if _, err := fmt.Sscanf(seg, "%d", &dayNum); err != nil {
			return nil, err
		}
		result = append(result, dayNum)
	}
	return result, nil
}
