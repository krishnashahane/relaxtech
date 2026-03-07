package cmd

import (
	"context"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/krishna/relaxtech/internal/client"
	"github.com/krishna/relaxtech/internal/output"
)

var deviceStatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Retrieve and display the current device state",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := ensureCredentials(); err != nil {
			return err
		}
		cl := client.New(viper.GetString("email"), viper.GetString("password"), viper.GetString("user_id"), viper.GetString("client_id"), viper.GetString("client_secret"))
		st, err := cl.GetStatus(context.Background())
		if err != nil {
			return err
		}
		row := map[string]any{"mode": st.CurrentState.Type, "level": st.CurrentLevel}
		columns := viper.GetStringSlice("fields")
		filtered := output.SelectColumns([]map[string]any{row}, columns)
		headers := columns
		if len(headers) == 0 {
			headers = []string{"mode", "level"}
		}
		return output.Render(output.DisplayMode(viper.GetString("output")), headers, filtered)
	},
}
