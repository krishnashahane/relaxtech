package cmd

import (
	"context"
	"fmt"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/krishna/relaxtech/internal/client"
	"github.com/krishna/relaxtech/internal/output"
)

// sleepRangeCmd fetches sleep metrics across a range of dates.
var sleepRangeCmd = &cobra.Command{
	Use:   "range",
	Short: "Fetch sleep metrics for a date range",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := ensureCredentials(); err != nil {
			return err
		}
		beginDate := viper.GetString("from")
		endDate := viper.GetString("to")
		if beginDate == "" || endDate == "" {
			return fmt.Errorf("--from and --to are required")
		}
		dateLayout := "2006-01-02"
		startDay, err := time.Parse(dateLayout, beginDate)
		if err != nil {
			return err
		}
		lastDay, err := time.Parse(dateLayout, endDate)
		if err != nil {
			return err
		}
		if lastDay.Before(startDay) {
			return fmt.Errorf("to must be >= from")
		}
		tzLabel := viper.GetString("timezone")
		if tzLabel == "local" {
			tzLabel = time.Local.String()
		}
		cl := client.New(viper.GetString("email"), viper.GetString("password"), viper.GetString("user_id"), viper.GetString("client_id"), viper.GetString("client_secret"))
		records := []map[string]any{}
		for cursor := startDay; !cursor.After(lastDay); cursor = cursor.Add(24 * time.Hour) {
			dayData, err := cl.GetSleepDay(context.Background(), cursor.Format(dateLayout), tzLabel)
			if err != nil {
				return err
			}
			records = append(records, map[string]any{
				"date":       dayData.Date,
				"score":      dayData.Score,
				"duration":   dayData.Duration,
				"tnt":        dayData.Tnt,
				"resp_rate":  dayData.Respiratory,
				"heart_rate": dayData.HeartRate,
				"hrv_score":  dayData.SleepQuality.HRV.Score,
			})
		}
		records = output.SelectColumns(records, viper.GetStringSlice("fields"))
		columnNames := []string{"date", "score", "duration", "tnt", "resp_rate", "heart_rate", "hrv_score"}
		if len(viper.GetStringSlice("fields")) > 0 {
			columnNames = viper.GetStringSlice("fields")
		}
		return output.Render(output.DisplayMode(viper.GetString("output")), columnNames, records)
	},
}

func init() {
	sleepRangeCmd.Flags().String("from", "", "start date YYYY-MM-DD")
	sleepRangeCmd.Flags().String("to", "", "end date YYYY-MM-DD")
	viper.BindPFlag("from", sleepRangeCmd.Flags().Lookup("from"))
	viper.BindPFlag("to", sleepRangeCmd.Flags().Lookup("to"))
	if sleepCmd != nil {
		sleepCmd.AddCommand(sleepRangeCmd)
	}
}
