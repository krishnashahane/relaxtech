package cmd

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/krishna/relaxtech/internal/client"
	"github.com/krishna/relaxtech/internal/daemon"
)

var temperatureCmd = &cobra.Command{
	Use:   "temp <value>",
	Short: "Adjust the pod temperature (e.g., 68F, 20C, or a raw level from -100 to 100)",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := ensureCredentials(); err != nil {
			return err
		}
		level, err := daemon.ConvertTemperature(args[0])
		if err != nil {
			return err
		}
		cl := client.New(viper.GetString("email"), viper.GetString("password"), viper.GetString("user_id"), viper.GetString("client_id"), viper.GetString("client_secret"))
		if err := cl.SetTemperature(context.Background(), level); err != nil {
			return err
		}
		fmt.Printf("temperature set (level %d)\n", level)
		return nil
	},
}
