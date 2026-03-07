package cmd

import (
	"context"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/krishna/relaxtech/internal/client"
	"github.com/krishna/relaxtech/internal/output"
)

// householdCmd groups household-related information subcommands.
var householdCmd = &cobra.Command{Use: "household", Short: "Household info"}

// buildHouseholdCmd creates a simple household subcommand that fetches data and renders it.
func buildHouseholdCmd(label string, fetchFn func(*client.Client, context.Context) (any, error)) *cobra.Command {
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
	householdCmd.AddCommand(
		buildHouseholdCmd("summary", func(cl *client.Client, ctx context.Context) (any, error) { return cl.Household().Summary(ctx) }),
		buildHouseholdCmd("schedule", func(cl *client.Client, ctx context.Context) (any, error) { return cl.Household().Schedule(ctx) }),
		buildHouseholdCmd("current-set", func(cl *client.Client, ctx context.Context) (any, error) { return cl.Household().CurrentSet(ctx) }),
		buildHouseholdCmd("invitations", func(cl *client.Client, ctx context.Context) (any, error) { return cl.Household().Invitations(ctx) }),
		buildHouseholdCmd("devices", func(cl *client.Client, ctx context.Context) (any, error) { return cl.Household().Devices(ctx) }),
		buildHouseholdCmd("users", func(cl *client.Client, ctx context.Context) (any, error) { return cl.Household().Users(ctx) }),
		buildHouseholdCmd("guests", func(cl *client.Client, ctx context.Context) (any, error) { return cl.Household().Guests(ctx) }),
	)
}
