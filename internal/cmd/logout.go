package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/krishna/relaxtech/internal/client"
	"github.com/krishna/relaxtech/internal/tokencache"
)

var signoutCmd = &cobra.Command{
	Use:   "logout",
	Short: "Remove stored authentication token",
	RunE: func(cmd *cobra.Command, args []string) error {
		cl := client.New(
			viper.GetString("email"),
			viper.GetString("password"),
			viper.GetString("user_id"),
			viper.GetString("client_id"),
			viper.GetString("client_secret"),
		)
		if err := tokencache.Remove(cl.AuthContext()); err != nil {
			return fmt.Errorf("remove stored token: %w", err)
		}
		fmt.Println("Logged out (stored token removed)")
		return nil
	},
}
