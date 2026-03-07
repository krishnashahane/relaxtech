package cmd

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/krishna/relaxtech/internal/client"
)

var activateCmd = &cobra.Command{
	Use:   "on",
	Short: "Power on the pod",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := ensureCredentials(); err != nil {
			return err
		}
		cl := client.New(viper.GetString("email"), viper.GetString("password"), viper.GetString("user_id"), viper.GetString("client_id"), viper.GetString("client_secret"))
		if err := cl.TurnOn(context.Background()); err != nil {
			return err
		}
		fmt.Println("pod turned on")
		return nil
	},
}
