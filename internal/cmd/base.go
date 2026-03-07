package cmd

import (
	"context"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/krishna/relaxtech/internal/client"
	"github.com/krishna/relaxtech/internal/output"
)

// baseCmd groups adjustable base control subcommands.
var baseCmd = &cobra.Command{Use: "base", Short: "Adjustable base controls"}

// baseInfoCmd retrieves information about the adjustable base.
var baseInfoCmd = &cobra.Command{Use: "info", RunE: func(cmd *cobra.Command, args []string) error {
	if err := ensureCredentials(); err != nil {
		return err
	}
	cl := client.New(viper.GetString("email"), viper.GetString("password"), viper.GetString("user_id"), viper.GetString("client_id"), viper.GetString("client_secret"))
	details, err := cl.Base().Info(context.Background())
	if err != nil {
		return err
	}
	return output.Render(output.DisplayMode(viper.GetString("output")), []string{"info"}, []map[string]any{{"info": details}})
}}

// baseAngleCmd sets the head and foot angles of the adjustable base.
var baseAngleCmd = &cobra.Command{Use: "angle", RunE: func(cmd *cobra.Command, args []string) error {
	if err := ensureCredentials(); err != nil {
		return err
	}
	headAngle := viper.GetInt("head")
	footAngle := viper.GetInt("foot")
	cl := client.New(viper.GetString("email"), viper.GetString("password"), viper.GetString("user_id"), viper.GetString("client_id"), viper.GetString("client_secret"))
	return cl.Base().SetAngle(context.Background(), headAngle, footAngle)
}}

// basePresetsCmd lists available base position presets.
var basePresetsCmd = &cobra.Command{Use: "presets", RunE: func(cmd *cobra.Command, args []string) error {
	if err := ensureCredentials(); err != nil {
		return err
	}
	cl := client.New(viper.GetString("email"), viper.GetString("password"), viper.GetString("user_id"), viper.GetString("client_id"), viper.GetString("client_secret"))
	presetList, err := cl.Base().Presets(context.Background())
	if err != nil {
		return err
	}
	return output.Render(output.DisplayMode(viper.GetString("output")), []string{"presets"}, []map[string]any{{"presets": presetList}})
}}

// basePresetRunCmd activates a named preset on the adjustable base.
var basePresetRunCmd = &cobra.Command{Use: "preset-run", RunE: func(cmd *cobra.Command, args []string) error {
	if err := ensureCredentials(); err != nil {
		return err
	}
	presetName := viper.GetString("name")
	cl := client.New(viper.GetString("email"), viper.GetString("password"), viper.GetString("user_id"), viper.GetString("client_id"), viper.GetString("client_secret"))
	return cl.Base().RunPreset(context.Background(), presetName)
}}

// baseTestCmd runs a vibration test on the adjustable base.
var baseTestCmd = &cobra.Command{Use: "test", RunE: func(cmd *cobra.Command, args []string) error {
	if err := ensureCredentials(); err != nil {
		return err
	}
	cl := client.New(viper.GetString("email"), viper.GetString("password"), viper.GetString("user_id"), viper.GetString("client_id"), viper.GetString("client_secret"))
	return cl.Base().VibrationTest(context.Background())
}}

func init() {
	baseAngleCmd.Flags().Int("head", 0, "head angle")
	baseAngleCmd.Flags().Int("foot", 0, "foot angle")
	viper.BindPFlag("head", baseAngleCmd.Flags().Lookup("head"))
	viper.BindPFlag("foot", baseAngleCmd.Flags().Lookup("foot"))
	basePresetRunCmd.Flags().String("name", "", "preset name")
	viper.BindPFlag("name", basePresetRunCmd.Flags().Lookup("name"))

	baseCmd.AddCommand(baseInfoCmd, baseAngleCmd, basePresetsCmd, basePresetRunCmd, baseTestCmd)
}
