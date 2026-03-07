package cmd

import (
	"context"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/krishna/relaxtech/internal/client"
	"github.com/krishna/relaxtech/internal/output"
)

// featsCmd retrieves and displays available release features.
var featsCmd = &cobra.Command{
	Use:   "feats",
	Short: "List release features",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := ensureCredentials(); err != nil {
			return err
		}
		cl := client.New(viper.GetString("email"), viper.GetString("password"), viper.GetString("user_id"), viper.GetString("client_id"), viper.GetString("client_secret"))
		featureList, err := cl.ReleaseFeatures(context.Background())
		if err != nil {
			return err
		}
		records := make([]map[string]any, 0, len(featureList))
		for _, feat := range featureList {
			records = append(records, map[string]any{"title": feat.Title, "body": feat.Body})
		}
		records = output.SelectColumns(records, viper.GetStringSlice("fields"))
		return output.Render(output.DisplayMode(viper.GetString("output")), []string{"title", "body"}, records)
	},
}
