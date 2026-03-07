package cmd

import (
	"context"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/krishna/relaxtech/internal/client"
	"github.com/krishna/relaxtech/internal/output"
)

// autopilotCmd groups autopilot configuration and history subcommands.
var autopilotCmd = &cobra.Command{Use: "autopilot", Short: "Autopilot settings"}

var (
	// autopilotDetailsCmd retrieves current autopilot configuration details.
	autopilotDetailsCmd = buildAutopilotCmd("details", func(cl *client.Client, ctx context.Context) (any, error) { return cl.Autopilot().Details(ctx) })
	// autopilotHistoryCmd retrieves autopilot adjustment history.
	autopilotHistoryCmd = buildAutopilotCmd("history", func(cl *client.Client, ctx context.Context) (any, error) { return cl.Autopilot().History(ctx) })
	// autopilotRecapCmd retrieves the autopilot performance recap.
	autopilotRecapCmd = buildAutopilotCmd("recap", func(cl *client.Client, ctx context.Context) (any, error) { return cl.Autopilot().Recap(ctx) })
)

// autopilotLevelCmd toggles automatic level suggestions on or off.
var autopilotLevelCmd = &cobra.Command{Use: "level-suggestions", RunE: func(cmd *cobra.Command, args []string) error {
	if err := ensureCredentials(); err != nil {
		return err
	}
	toggle := viper.GetBool("enabled")
	cl := client.New(viper.GetString("email"), viper.GetString("password"), viper.GetString("user_id"), viper.GetString("client_id"), viper.GetString("client_secret"))
	return cl.Autopilot().SetLevelSuggestions(context.Background(), toggle)
}}

// autopilotSnoreCmd toggles snore mitigation on or off.
var autopilotSnoreCmd = &cobra.Command{Use: "snore-mitigation", RunE: func(cmd *cobra.Command, args []string) error {
	if err := ensureCredentials(); err != nil {
		return err
	}
	toggle := viper.GetBool("enabled")
	cl := client.New(viper.GetString("email"), viper.GetString("password"), viper.GetString("user_id"), viper.GetString("client_id"), viper.GetString("client_secret"))
	return cl.Autopilot().SetSnoreMitigation(context.Background(), toggle)
}}

// buildAutopilotCmd creates a simple cobra command that fetches autopilot data and renders it.
func buildAutopilotCmd(label string, fetchFn func(*client.Client, context.Context) (any, error)) *cobra.Command {
	return &cobra.Command{Use: label, RunE: func(cmd *cobra.Command, args []string) error {
		if err := ensureCredentials(); err != nil {
			return err
		}
		cl := client.New(viper.GetString("email"), viper.GetString("password"), viper.GetString("user_id"), viper.GetString("client_id"), viper.GetString("client_secret"))
		payload, err := fetchFn(cl, context.Background())
		if err != nil {
			return err
		}
		return output.Render(output.DisplayMode(viper.GetString("output")), []string{label}, []map[string]any{{label: payload}})
	}}
}

func init() {
	autopilotLevelCmd.Flags().Bool("enabled", true, "enable or disable")
	viper.BindPFlag("enabled", autopilotLevelCmd.Flags().Lookup("enabled"))
	autopilotSnoreCmd.Flags().Bool("enabled", true, "enable or disable")
	viper.BindPFlag("enabled", autopilotSnoreCmd.Flags().Lookup("enabled"))

	autopilotCmd.AddCommand(autopilotDetailsCmd, autopilotHistoryCmd, autopilotRecapCmd, autopilotLevelCmd, autopilotSnoreCmd)
}
