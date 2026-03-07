package cmd

import (
	"context"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/krishna/relaxtech/internal/client"
	"github.com/krishna/relaxtech/internal/output"
)

// metricsCmd groups sleep metrics and insight subcommands.
var metricsCmd = &cobra.Command{Use: "metrics", Short: "Sleep metrics and insights"}

// metricsTrendsCmd fetches sleep trend data over a date range.
var metricsTrendsCmd = &cobra.Command{Use: "trends", RunE: func(cmd *cobra.Command, args []string) error {
	if err := ensureCredentials(); err != nil {
		return err
	}
	startDate := viper.GetString("from")
	endDate := viper.GetString("to")
	cl := client.New(viper.GetString("email"), viper.GetString("password"), viper.GetString("user_id"), viper.GetString("client_id"), viper.GetString("client_secret"))
	var trendData any
	if err := cl.Metrics().Trends(context.Background(), startDate, endDate, &trendData); err != nil {
		return err
	}
	return output.Render(output.DisplayMode(viper.GetString("output")), []string{"trends"}, []map[string]any{{"trends": trendData}})
}}

// metricsIntervalsCmd fetches interval-level data for a specific session.
var metricsIntervalsCmd = &cobra.Command{Use: "intervals", RunE: func(cmd *cobra.Command, args []string) error {
	if err := ensureCredentials(); err != nil {
		return err
	}
	sessionID := viper.GetString("id")
	cl := client.New(viper.GetString("email"), viper.GetString("password"), viper.GetString("user_id"), viper.GetString("client_id"), viper.GetString("client_secret"))
	var intervalData any
	if err := cl.Metrics().Intervals(context.Background(), sessionID, &intervalData); err != nil {
		return err
	}
	return output.Render(output.DisplayMode(viper.GetString("output")), []string{"interval"}, []map[string]any{{"interval": intervalData}})
}}

// metricsSummaryCmd retrieves a high-level sleep metrics summary.
var metricsSummaryCmd = &cobra.Command{Use: "summary", RunE: func(cmd *cobra.Command, args []string) error {
	if err := ensureCredentials(); err != nil {
		return err
	}
	cl := client.New(viper.GetString("email"), viper.GetString("password"), viper.GetString("user_id"), viper.GetString("client_id"), viper.GetString("client_secret"))
	var summaryData any
	if err := cl.Metrics().Summary(context.Background(), &summaryData); err != nil {
		return err
	}
	return output.Render(output.DisplayMode(viper.GetString("output")), []string{"summary"}, []map[string]any{{"summary": summaryData}})
}}

// metricsAggregateCmd retrieves aggregated sleep metrics.
var metricsAggregateCmd = &cobra.Command{Use: "aggregate", RunE: func(cmd *cobra.Command, args []string) error {
	if err := ensureCredentials(); err != nil {
		return err
	}
	cl := client.New(viper.GetString("email"), viper.GetString("password"), viper.GetString("user_id"), viper.GetString("client_id"), viper.GetString("client_secret"))
	var aggregateData any
	if err := cl.Metrics().Aggregate(context.Background(), &aggregateData); err != nil {
		return err
	}
	return output.Render(output.DisplayMode(viper.GetString("output")), []string{"aggregate"}, []map[string]any{{"aggregate": aggregateData}})
}}

// metricsInsightsCmd retrieves personalized sleep insights.
var metricsInsightsCmd = &cobra.Command{Use: "insights", RunE: func(cmd *cobra.Command, args []string) error {
	if err := ensureCredentials(); err != nil {
		return err
	}
	cl := client.New(viper.GetString("email"), viper.GetString("password"), viper.GetString("user_id"), viper.GetString("client_id"), viper.GetString("client_secret"))
	var insightData any
	if err := cl.Metrics().Insights(context.Background(), &insightData); err != nil {
		return err
	}
	return output.Render(output.DisplayMode(viper.GetString("output")), []string{"insights"}, []map[string]any{{"insights": insightData}})
}}

func init() {
	metricsTrendsCmd.Flags().String("from", "", "from date YYYY-MM-DD")
	metricsTrendsCmd.Flags().String("to", "", "to date YYYY-MM-DD")
	viper.BindPFlag("from", metricsTrendsCmd.Flags().Lookup("from"))
	viper.BindPFlag("to", metricsTrendsCmd.Flags().Lookup("to"))
	metricsIntervalsCmd.Flags().String("id", "", "session id")
	viper.BindPFlag("id", metricsIntervalsCmd.Flags().Lookup("id"))

	metricsCmd.AddCommand(metricsTrendsCmd, metricsIntervalsCmd, metricsSummaryCmd, metricsAggregateCmd, metricsInsightsCmd)
}
