package cmd

import (
	"context"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/krishna/relaxtech/internal/client"
	"github.com/krishna/relaxtech/internal/output"
)

// sleepCmd groups sleep analytics subcommands.
var sleepCmd = &cobra.Command{
	Use:   "sleep",
	Short: "Sleep analytics commands",
}

// sleepDayCmd fetches sleep metrics for a single date.
var sleepDayCmd = &cobra.Command{
	Use:   "day",
	Short: "Fetch sleep metrics for a date (YYYY-MM-DD)",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := ensureCredentials(); err != nil {
			return err
		}
		targetDate := viper.GetString("date")
		if targetDate == "" {
			targetDate = time.Now().Format("2006-01-02")
		}
		tzLabel := viper.GetString("timezone")
		if tzLabel == "local" {
			tzLabel = time.Local.String()
		}
		cl := client.New(viper.GetString("email"), viper.GetString("password"), viper.GetString("user_id"), viper.GetString("client_id"), viper.GetString("client_secret"))
		dayData, err := cl.GetSleepDay(context.Background(), targetDate, tzLabel)
		if err != nil {
			return err
		}
		records := []map[string]any{
			{
				"date":           dayData.Date,
				"score":          dayData.Score,
				"tnt":            dayData.Tnt,
				"resp_rate":      dayData.Respiratory,
				"heart_rate":     dayData.HeartRate,
				"duration":       dayData.Duration,
				"latency_asleep": dayData.LatencyAsleep,
				"latency_out":    dayData.LatencyOut,
				"hrv_score":      dayData.SleepQuality.HRV.Score,
			},
		}
		records = output.SelectColumns(records, viper.GetStringSlice("fields"))
		return output.Render(output.DisplayMode(viper.GetString("output")), []string{"date", "score", "duration", "latency_asleep", "latency_out", "tnt", "resp_rate", "heart_rate", "hrv_score"}, records)
	},
}

func init() {
	sleepCmd.PersistentFlags().String("date", "", "date YYYY-MM-DD (default today)")
	viper.BindPFlag("date", sleepCmd.PersistentFlags().Lookup("date"))
	sleepCmd.AddCommand(sleepDayCmd)
}
