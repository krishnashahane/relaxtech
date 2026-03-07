package cmd

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/krishna/relaxtech/internal/client"
)

// presenceCmd checks whether the user is currently detected in bed.
var presenceCmd = &cobra.Command{
	Use:   "presence",
	Short: "Check if user is in bed",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := ensureCredentials(); err != nil {
			return err
		}
		cl := client.New(viper.GetString("email"), viper.GetString("password"), viper.GetString("user_id"), viper.GetString("client_id"), viper.GetString("client_secret"))
		detected, err := cl.GetPresence(context.Background())
		if err != nil {
			return err
		}
		fmt.Printf("present: %v\n", detected)
		return nil
	},
}
