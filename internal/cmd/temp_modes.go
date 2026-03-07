package cmd

import (
	"context"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/krishna/relaxtech/internal/client"
	"github.com/krishna/relaxtech/internal/output"
)

// tempModeCmd groups temperature mode subcommands (nap, hot-flash, events).
var tempModeCmd = &cobra.Command{
	Use:   "tempmode",
	Short: "Temperature modes (nap, hot-flash, temp events)",
}

var (
	// tempNapCmd groups nap mode control subcommands.
	tempNapCmd = &cobra.Command{Use: "nap", Short: "Nap mode controls"}

	// tempNapOnCmd activates nap mode on the device.
	tempNapOnCmd = &cobra.Command{Use: "on", RunE: func(cmd *cobra.Command, args []string) error {
		if err := ensureCredentials(); err != nil {
			return err
		}
		cl := client.New(viper.GetString("email"), viper.GetString("password"), viper.GetString("user_id"), viper.GetString("client_id"), viper.GetString("client_secret"))
		return cl.TempModes().NapActivate(context.Background())
	}}
)

// tempNapOffCmd deactivates nap mode on the device.
var tempNapOffCmd = &cobra.Command{Use: "off", RunE: func(cmd *cobra.Command, args []string) error {
	if err := ensureCredentials(); err != nil {
		return err
	}
	cl := client.New(viper.GetString("email"), viper.GetString("password"), viper.GetString("user_id"), viper.GetString("client_id"), viper.GetString("client_secret"))
	return cl.TempModes().NapDeactivate(context.Background())
}}

// tempNapExtendCmd extends the current nap session.
var tempNapExtendCmd = &cobra.Command{Use: "extend", RunE: func(cmd *cobra.Command, args []string) error {
	if err := ensureCredentials(); err != nil {
		return err
	}
	cl := client.New(viper.GetString("email"), viper.GetString("password"), viper.GetString("user_id"), viper.GetString("client_id"), viper.GetString("client_secret"))
	return cl.TempModes().NapExtend(context.Background())
}}

// tempNapStatusCmd displays the current nap mode status.
var tempNapStatusCmd = &cobra.Command{Use: "status", RunE: func(cmd *cobra.Command, args []string) error {
	if err := ensureCredentials(); err != nil {
		return err
	}
	cl := client.New(viper.GetString("email"), viper.GetString("password"), viper.GetString("user_id"), viper.GetString("client_id"), viper.GetString("client_secret"))
	var statusData map[string]any
	if err := cl.TempModes().NapStatus(context.Background(), &statusData); err != nil {
		return err
	}
	records := output.SelectColumns([]map[string]any{statusData}, viper.GetStringSlice("fields"))
	columnNames := viper.GetStringSlice("fields")
	if len(columnNames) == 0 {
		columnNames = sortedMapKeys(statusData)
	}
	return output.Render(output.DisplayMode(viper.GetString("output")), columnNames, records)
}}

var (
	// tempHotCmd groups hot-flash mode control subcommands.
	tempHotCmd = &cobra.Command{Use: "hotflash", Short: "Hot-flash mode controls"}

	// tempHotOnCmd activates hot-flash mode.
	tempHotOnCmd = &cobra.Command{Use: "on", RunE: func(cmd *cobra.Command, args []string) error {
		if err := ensureCredentials(); err != nil {
			return err
		}
		cl := client.New(viper.GetString("email"), viper.GetString("password"), viper.GetString("user_id"), viper.GetString("client_id"), viper.GetString("client_secret"))
		return cl.TempModes().HotFlashActivate(context.Background())
	}}
)

// tempHotOffCmd deactivates hot-flash mode.
var tempHotOffCmd = &cobra.Command{Use: "off", RunE: func(cmd *cobra.Command, args []string) error {
	if err := ensureCredentials(); err != nil {
		return err
	}
	cl := client.New(viper.GetString("email"), viper.GetString("password"), viper.GetString("user_id"), viper.GetString("client_id"), viper.GetString("client_secret"))
	return cl.TempModes().HotFlashDeactivate(context.Background())
}}

// tempHotStatusCmd displays the current hot-flash mode status.
var tempHotStatusCmd = &cobra.Command{Use: "status", RunE: func(cmd *cobra.Command, args []string) error {
	if err := ensureCredentials(); err != nil {
		return err
	}
	cl := client.New(viper.GetString("email"), viper.GetString("password"), viper.GetString("user_id"), viper.GetString("client_id"), viper.GetString("client_secret"))
	var statusData map[string]any
	if err := cl.TempModes().HotFlashStatus(context.Background(), &statusData); err != nil {
		return err
	}
	records := output.SelectColumns([]map[string]any{statusData}, viper.GetStringSlice("fields"))
	columnNames := viper.GetStringSlice("fields")
	if len(columnNames) == 0 {
		columnNames = sortedMapKeys(statusData)
	}
	return output.Render(output.DisplayMode(viper.GetString("output")), columnNames, records)
}}

// tempEventsCmd retrieves temperature events within an optional date range.
var tempEventsCmd = &cobra.Command{Use: "events", RunE: func(cmd *cobra.Command, args []string) error {
	if err := ensureCredentials(); err != nil {
		return err
	}
	cl := client.New(viper.GetString("email"), viper.GetString("password"), viper.GetString("user_id"), viper.GetString("client_id"), viper.GetString("client_secret"))
	startDate := viper.GetString("from")
	endDate := viper.GetString("to")
	var eventData any
	if err := cl.TempModes().TempEvents(context.Background(), startDate, endDate, &eventData); err != nil {
		return err
	}
	records := []map[string]any{{"events": eventData}}
	records = output.SelectColumns(records, viper.GetStringSlice("fields"))
	columnNames := viper.GetStringSlice("fields")
	if len(columnNames) == 0 {
		columnNames = sortedMapKeys(records[0])
	}
	return output.Render(output.DisplayMode(viper.GetString("output")), columnNames, records)
}}

func init() {
	tempNapCmd.AddCommand(tempNapOnCmd, tempNapOffCmd, tempNapExtendCmd, tempNapStatusCmd)
	tempHotCmd.AddCommand(tempHotOnCmd, tempHotOffCmd, tempHotStatusCmd)
	tempEventsCmd.Flags().String("from", "", "from date (YYYY-MM-DD)")
	tempEventsCmd.Flags().String("to", "", "to date (YYYY-MM-DD)")
	viper.BindPFlag("from", tempEventsCmd.Flags().Lookup("from"))
	viper.BindPFlag("to", tempEventsCmd.Flags().Lookup("to"))

	tempModeCmd.AddCommand(tempNapCmd, tempHotCmd, tempEventsCmd)
}
