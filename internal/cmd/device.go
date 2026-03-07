package cmd

import (
	"context"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/krishna/relaxtech/internal/client"
	"github.com/krishna/relaxtech/internal/output"
)

// deviceCmd groups device information and priming subcommands.
var deviceCmd = &cobra.Command{Use: "device", Short: "Device info and priming"}

// buildDeviceCmd creates a simple device subcommand that fetches data and renders it.
func buildDeviceCmd(label string, fetchFn func(ctx context.Context) (any, error)) *cobra.Command {
	return &cobra.Command{Use: label, RunE: func(cmd *cobra.Command, args []string) error {
		if err := ensureCredentials(); err != nil {
			return err
		}
		payload, err := fetchFn(cmd.Context())
		if err != nil {
			return err
		}
		return output.Render(output.DisplayMode(viper.GetString("output")), []string{label}, []map[string]any{{label: payload}})
	}}
}

func init() {
	deviceCmd.AddCommand(
		buildDeviceCmd("info", func(ctx context.Context) (any, error) {
			cl := client.New(viper.GetString("email"), viper.GetString("password"), viper.GetString("user_id"), viper.GetString("client_id"), viper.GetString("client_secret"))
			return cl.Device().Info(ctx)
		}),
		buildDeviceCmd("peripherals", func(ctx context.Context) (any, error) {
			cl := client.New(viper.GetString("email"), viper.GetString("password"), viper.GetString("user_id"), viper.GetString("client_id"), viper.GetString("client_secret"))
			return cl.Device().Peripherals(ctx)
		}),
		buildDeviceCmd("owner", func(ctx context.Context) (any, error) {
			cl := client.New(viper.GetString("email"), viper.GetString("password"), viper.GetString("user_id"), viper.GetString("client_id"), viper.GetString("client_secret"))
			return cl.Device().Owner(ctx)
		}),
		buildDeviceCmd("warranty", func(ctx context.Context) (any, error) {
			cl := client.New(viper.GetString("email"), viper.GetString("password"), viper.GetString("user_id"), viper.GetString("client_id"), viper.GetString("client_secret"))
			return cl.Device().Warranty(ctx)
		}),
		buildDeviceCmd("online", func(ctx context.Context) (any, error) {
			cl := client.New(viper.GetString("email"), viper.GetString("password"), viper.GetString("user_id"), viper.GetString("client_id"), viper.GetString("client_secret"))
			return cl.Device().Online(ctx)
		}),
		buildDeviceCmd("priming-tasks", func(ctx context.Context) (any, error) {
			cl := client.New(viper.GetString("email"), viper.GetString("password"), viper.GetString("user_id"), viper.GetString("client_id"), viper.GetString("client_secret"))
			return cl.Device().PrimingTasks(ctx)
		}),
		buildDeviceCmd("priming-schedule", func(ctx context.Context) (any, error) {
			cl := client.New(viper.GetString("email"), viper.GetString("password"), viper.GetString("user_id"), viper.GetString("client_id"), viper.GetString("client_secret"))
			return cl.Device().PrimingSchedule(ctx)
		}),
	)
}
