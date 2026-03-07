package cmd

import (
	"context"
	"fmt"
	"sort"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/krishna/relaxtech/internal/client"
	"github.com/krishna/relaxtech/internal/output"
)

// scheduleCmd groups cloud-based temperature schedule management subcommands.
var scheduleCmd = &cobra.Command{
	Use:   "schedule",
	Short: "Manage device temperature schedules (cloud)",
}

// scheduleNextCmd displays upcoming schedule events sorted by time.
var scheduleNextCmd = &cobra.Command{
	Use:   "next",
	Short: "Show next upcoming schedule events",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := ensureCredentials(); err != nil {
			return err
		}
		tzLabel := viper.GetString("timezone")
		region := time.Local
		if tzLabel != "" && tzLabel != "local" {
			parsed, err := time.LoadLocation(tzLabel)
			if err != nil {
				return err
			}
			region = parsed
		}
		currentTime := time.Now().In(region)

		cl := client.New(viper.GetString("email"), viper.GetString("password"), viper.GetString("user_id"), viper.GetString("client_id"), viper.GetString("client_secret"))
		schedEntries, err := cl.ListSchedules(context.Background())
		if err != nil {
			return err
		}

		records := make([]map[string]any, 0, len(schedEntries))
		for _, entry := range schedEntries {
			upcoming := computeNextOccurrence(currentTime, entry, region)
			records = append(records, map[string]any{
				"id":      entry.ID,
				"start":   entry.StartTime,
				"days":    entry.DaysOfWeek,
				"level":   entry.Level,
				"enabled": entry.Enabled,
				"next":    upcoming.Format(time.RFC3339),
			})
		}

		sort.Slice(records, func(i, j int) bool { return records[i]["next"].(string) < records[j]["next"].(string) })
		records = output.SelectColumns(records, viper.GetStringSlice("fields"))
		columnNames := []string{"id", "start", "days", "level", "enabled", "next"}
		if len(viper.GetStringSlice("fields")) > 0 {
			columnNames = viper.GetStringSlice("fields")
		}
		return output.Render(output.DisplayMode(viper.GetString("output")), columnNames, records)
	},
}

// scheduleListCmd retrieves and displays all configured schedules.
var scheduleListCmd = &cobra.Command{
	Use:   "list",
	Short: "List schedules",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := ensureCredentials(); err != nil {
			return err
		}
		cl := client.New(viper.GetString("email"), viper.GetString("password"), viper.GetString("user_id"), viper.GetString("client_id"), viper.GetString("client_secret"))
		schedEntries, err := cl.ListSchedules(context.Background())
		if err != nil {
			return err
		}
		records := make([]map[string]any, 0, len(schedEntries))
		for _, entry := range schedEntries {
			records = append(records, map[string]any{
				"id":      entry.ID,
				"start":   entry.StartTime,
				"level":   entry.Level,
				"days":    entry.DaysOfWeek,
				"enabled": entry.Enabled,
			})
		}
		records = output.SelectColumns(records, viper.GetStringSlice("fields"))
		return output.Render(output.DisplayMode(viper.GetString("output")), []string{"id", "start", "level", "days", "enabled"}, records)
	},
}

// scheduleCreateCmd adds a new temperature schedule entry.
var scheduleCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create schedule",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := ensureCredentials(); err != nil {
			return err
		}
		startVal := viper.GetString("start")
		if startVal == "" {
			return fmt.Errorf("--start HH:MM required")
		}
		tempLevel := viper.GetInt("level")
		weekdays := viper.GetIntSlice("days")
		if len(weekdays) == 0 {
			return fmt.Errorf("--days required")
		}
		active := !viper.GetBool("disabled")
		cl := client.New(viper.GetString("email"), viper.GetString("password"), viper.GetString("user_id"), viper.GetString("client_id"), viper.GetString("client_secret"))
		entry := client.TemperatureSchedule{StartTime: startVal, Level: tempLevel, DaysOfWeek: weekdays, Enabled: active}
		created, err := cl.CreateSchedule(context.Background(), entry)
		if err != nil {
			return err
		}
		fmt.Printf("created schedule %s\n", created.ID)
		return nil
	},
}

// scheduleUpdateCmd modifies an existing schedule entry by its ID.
var scheduleUpdateCmd = &cobra.Command{
	Use:   "update <id>",
	Short: "Update schedule",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := ensureCredentials(); err != nil {
			return err
		}
		changes := map[string]any{}
		if cmd.Flags().Changed("start") {
			changes["startTime"] = viper.GetString("start")
		}
		if cmd.Flags().Changed("level") {
			changes["level"] = viper.GetInt("level")
		}
		if cmd.Flags().Changed("days") {
			changes["daysOfWeek"] = viper.GetIntSlice("days")
		}
		if cmd.Flags().Changed("enabled") {
			changes["enabled"] = viper.GetBool("enabled")
		}
		if len(changes) == 0 {
			return fmt.Errorf("no fields to update")
		}
		cl := client.New(viper.GetString("email"), viper.GetString("password"), viper.GetString("user_id"), viper.GetString("client_id"), viper.GetString("client_secret"))
		if _, err := cl.UpdateSchedule(context.Background(), args[0], changes); err != nil {
			return err
		}
		fmt.Println("updated")
		return nil
	},
}

// scheduleDeleteCmd removes a schedule entry by its ID.
var scheduleDeleteCmd = &cobra.Command{
	Use:   "delete <id>",
	Short: "Delete schedule",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := ensureCredentials(); err != nil {
			return err
		}
		cl := client.New(viper.GetString("email"), viper.GetString("password"), viper.GetString("user_id"), viper.GetString("client_id"), viper.GetString("client_secret"))
		if err := cl.DeleteSchedule(context.Background(), args[0]); err != nil {
			return err
		}
		fmt.Println("deleted")
		return nil
	},
}

func init() {
	scheduleCreateCmd.Flags().String("start", "", "HH:MM start time")
	scheduleCreateCmd.Flags().Int("level", 0, "Temperature level -100..100")
	scheduleCreateCmd.Flags().IntSlice("days", nil, "Comma-separated days 0=Sun..6=Sat")
	scheduleCreateCmd.Flags().Bool("disabled", false, "Create disabled")
	viper.BindPFlag("start", scheduleCreateCmd.Flags().Lookup("start"))
	viper.BindPFlag("level", scheduleCreateCmd.Flags().Lookup("level"))
	viper.BindPFlag("days", scheduleCreateCmd.Flags().Lookup("days"))
	viper.BindPFlag("disabled", scheduleCreateCmd.Flags().Lookup("disabled"))

	scheduleUpdateCmd.Flags().String("start", "", "HH:MM start time")
	scheduleUpdateCmd.Flags().Int("level", 0, "Temperature level -100..100")
	scheduleUpdateCmd.Flags().IntSlice("days", nil, "Comma-separated days 0=Sun..6=Sat")
	scheduleUpdateCmd.Flags().Bool("enabled", true, "Enable/disable schedule")
	viper.BindPFlag("start", scheduleUpdateCmd.Flags().Lookup("start"))
	viper.BindPFlag("level", scheduleUpdateCmd.Flags().Lookup("level"))
	viper.BindPFlag("days", scheduleUpdateCmd.Flags().Lookup("days"))
	viper.BindPFlag("enabled", scheduleUpdateCmd.Flags().Lookup("enabled"))

	scheduleCmd.AddCommand(scheduleListCmd, scheduleCreateCmd, scheduleUpdateCmd, scheduleDeleteCmd, scheduleNextCmd)
}

// computeNextOccurrence finds the next time a schedule entry will fire,
// scanning up to 14 days forward from the given reference time.
func computeNextOccurrence(refTime time.Time, entry client.TemperatureSchedule, region *time.Location) time.Time {
	schedHour, schedMin, _ := time.Now().Clock()
	if parsed, err := time.Parse("15:04", entry.StartTime); err == nil {
		schedHour, schedMin, _ = parsed.Clock()
	}
	activeDays := map[int]bool{}
	for _, d := range entry.DaysOfWeek {
		activeDays[d] = true
	}
	for offset := 0; offset < 14; offset++ {
		candidate := refTime.In(region).AddDate(0, 0, offset)
		if len(activeDays) > 0 && !activeDays[int(candidate.Weekday())] {
			continue
		}
		target := time.Date(candidate.Year(), candidate.Month(), candidate.Day(), schedHour, schedMin, 0, 0, region)
		if target.After(refTime) {
			return target
		}
	}
	return refTime
}
