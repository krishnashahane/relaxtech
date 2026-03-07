package cmd

import (
	"context"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/krishna/relaxtech/internal/client"
	"github.com/krishna/relaxtech/internal/output"
)

// tracksCmd lists all available audio tracks as a top-level command.
var tracksCmd = &cobra.Command{
	Use:   "tracks",
	Short: "List audio tracks",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := ensureCredentials(); err != nil {
			return err
		}
		cl := client.New(viper.GetString("email"), viper.GetString("password"), viper.GetString("user_id"), viper.GetString("client_id"), viper.GetString("client_secret"))
		trackList, err := cl.ListTracks(context.Background())
		if err != nil {
			return err
		}
		records := make([]map[string]any, 0, len(trackList))
		for _, t := range trackList {
			records = append(records, map[string]any{"id": t.ID, "title": t.Title, "type": t.Type})
		}
		selectedCols := viper.GetStringSlice("fields")
		records = output.SelectColumns(records, selectedCols)
		return output.Render(output.DisplayMode(viper.GetString("output")), []string{"id", "title", "type"}, records)
	},
}
