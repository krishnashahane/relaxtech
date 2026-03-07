package cmd

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/krishna/relaxtech/internal/client"
	"github.com/krishna/relaxtech/internal/output"
)

// audioCmd groups audio playback and track management subcommands.
var audioCmd = &cobra.Command{Use: "audio", Short: "Audio tracks and player"}

// audioTracksCmd lists all available audio tracks.
var audioTracksCmd = &cobra.Command{Use: "tracks", RunE: func(cmd *cobra.Command, args []string) error {
	if err := ensureCredentials(); err != nil {
		return err
	}
	cl := client.New(viper.GetString("email"), viper.GetString("password"), viper.GetString("user_id"), viper.GetString("client_id"), viper.GetString("client_secret"))
	trackList, err := cl.Audio().Tracks(context.Background())
	if err != nil {
		return err
	}
	records := make([]map[string]any, 0, len(trackList))
	for _, t := range trackList {
		records = append(records, map[string]any{"id": t.ID, "title": t.Title, "type": t.Type})
	}
	records = output.SelectColumns(records, viper.GetStringSlice("fields"))
	columnNames := viper.GetStringSlice("fields")
	if len(columnNames) == 0 {
		columnNames = []string{"id", "title", "type"}
	}
	return output.Render(output.DisplayMode(viper.GetString("output")), columnNames, records)
}}

// audioCategoriesCmd fetches the available audio categories.
var audioCategoriesCmd = &cobra.Command{Use: "categories", RunE: func(cmd *cobra.Command, args []string) error {
	if err := ensureCredentials(); err != nil {
		return err
	}
	cl := client.New(viper.GetString("email"), viper.GetString("password"), viper.GetString("user_id"), viper.GetString("client_id"), viper.GetString("client_secret"))
	categoryData, err := cl.Audio().Categories(context.Background())
	if err != nil {
		return err
	}
	records := []map[string]any{{"data": categoryData}}
	return output.Render(output.DisplayMode(viper.GetString("output")), []string{"data"}, records)
}}

// audioStateCmd shows the current player state.
var audioStateCmd = &cobra.Command{Use: "state", RunE: func(cmd *cobra.Command, args []string) error {
	if err := ensureCredentials(); err != nil {
		return err
	}
	cl := client.New(viper.GetString("email"), viper.GetString("password"), viper.GetString("user_id"), viper.GetString("client_id"), viper.GetString("client_secret"))
	playerState, err := cl.Audio().PlayerState(context.Background())
	if err != nil {
		return err
	}
	records := []map[string]any{{"state": playerState}}
	return output.Render(output.DisplayMode(viper.GetString("output")), []string{"state"}, records)
}}

// audioPlayCmd starts playback of a specified track.
var audioPlayCmd = &cobra.Command{Use: "play", RunE: func(cmd *cobra.Command, args []string) error {
	if err := ensureCredentials(); err != nil {
		return err
	}
	trackID := viper.GetString("track")
	cl := client.New(viper.GetString("email"), viper.GetString("password"), viper.GetString("user_id"), viper.GetString("client_id"), viper.GetString("client_secret"))
	return cl.Audio().Play(context.Background(), trackID)
}}

// audioPauseCmd pauses the currently playing audio.
var audioPauseCmd = &cobra.Command{Use: "pause", RunE: func(cmd *cobra.Command, args []string) error {
	if err := ensureCredentials(); err != nil {
		return err
	}
	cl := client.New(viper.GetString("email"), viper.GetString("password"), viper.GetString("user_id"), viper.GetString("client_id"), viper.GetString("client_secret"))
	return cl.Audio().Pause(context.Background())
}}

// audioSeekCmd jumps to a specific position in the current track.
var audioSeekCmd = &cobra.Command{Use: "seek", RunE: func(cmd *cobra.Command, args []string) error {
	if err := ensureCredentials(); err != nil {
		return err
	}
	positionMs := viper.GetInt("position")
	cl := client.New(viper.GetString("email"), viper.GetString("password"), viper.GetString("user_id"), viper.GetString("client_id"), viper.GetString("client_secret"))
	return cl.Audio().Seek(context.Background(), positionMs)
}}

// audioVolumeCmd adjusts the playback volume level.
var audioVolumeCmd = &cobra.Command{Use: "volume", RunE: func(cmd *cobra.Command, args []string) error {
	if err := ensureCredentials(); err != nil {
		return err
	}
	volumeLevel := viper.GetInt("level")
	cl := client.New(viper.GetString("email"), viper.GetString("password"), viper.GetString("user_id"), viper.GetString("client_id"), viper.GetString("client_secret"))
	return cl.Audio().Volume(context.Background(), volumeLevel)
}}

// audioPairCmd initiates Bluetooth pairing for the audio device.
var audioPairCmd = &cobra.Command{Use: "pair", RunE: func(cmd *cobra.Command, args []string) error {
	if err := ensureCredentials(); err != nil {
		return err
	}
	cl := client.New(viper.GetString("email"), viper.GetString("password"), viper.GetString("user_id"), viper.GetString("client_id"), viper.GetString("client_secret"))
	return cl.Audio().Pair(context.Background())
}}

// audioNextCmd fetches the recommended next track to play.
var audioNextCmd = &cobra.Command{Use: "next", RunE: func(cmd *cobra.Command, args []string) error {
	if err := ensureCredentials(); err != nil {
		return err
	}
	cl := client.New(viper.GetString("email"), viper.GetString("password"), viper.GetString("user_id"), viper.GetString("client_id"), viper.GetString("client_secret"))
	suggestion, err := cl.Audio().RecommendedNext(context.Background())
	if err != nil {
		return err
	}
	records := []map[string]any{{"next": suggestion}}
	return output.Render(output.DisplayMode(viper.GetString("output")), []string{"next"}, records)
}}

// audioFavoritesCmd groups favorite track management subcommands.
var audioFavoritesCmd = &cobra.Command{Use: "favorites", Short: "Favorite tracks"}

// audioFavListCmd lists all favorited tracks.
var audioFavListCmd = &cobra.Command{Use: "list", RunE: func(cmd *cobra.Command, args []string) error {
	if err := ensureCredentials(); err != nil {
		return err
	}
	cl := client.New(viper.GetString("email"), viper.GetString("password"), viper.GetString("user_id"), viper.GetString("client_id"), viper.GetString("client_secret"))
	favList, err := cl.Audio().Favorites(context.Background())
	if err != nil {
		return err
	}
	return output.Render(output.DisplayMode(viper.GetString("output")), []string{"favorites"}, []map[string]any{{"favorites": favList}})
}}

// audioFavAddCmd marks a track as a favorite.
var audioFavAddCmd = &cobra.Command{Use: "add", RunE: func(cmd *cobra.Command, args []string) error {
	if err := ensureCredentials(); err != nil {
		return err
	}
	trackID := viper.GetString("track")
	if trackID == "" {
		return fmt.Errorf("--track required")
	}
	cl := client.New(viper.GetString("email"), viper.GetString("password"), viper.GetString("user_id"), viper.GetString("client_id"), viper.GetString("client_secret"))
	return cl.Audio().AddFavorite(context.Background(), trackID)
}}

// audioFavRemoveCmd removes a track from favorites.
var audioFavRemoveCmd = &cobra.Command{Use: "remove", RunE: func(cmd *cobra.Command, args []string) error {
	if err := ensureCredentials(); err != nil {
		return err
	}
	trackID := viper.GetString("track")
	if trackID == "" {
		return fmt.Errorf("--track required")
	}
	cl := client.New(viper.GetString("email"), viper.GetString("password"), viper.GetString("user_id"), viper.GetString("client_id"), viper.GetString("client_secret"))
	return cl.Audio().RemoveFavorite(context.Background(), trackID)
}}

func init() {
	audioPlayCmd.Flags().String("track", "", "track ID to play")
	viper.BindPFlag("track", audioPlayCmd.Flags().Lookup("track"))
	audioSeekCmd.Flags().Int("position", 0, "position milliseconds")
	viper.BindPFlag("position", audioSeekCmd.Flags().Lookup("position"))
	audioVolumeCmd.Flags().Int("level", 50, "volume level 0-100")
	viper.BindPFlag("level", audioVolumeCmd.Flags().Lookup("level"))
	audioFavAddCmd.Flags().String("track", "", "track id")
	viper.BindPFlag("track", audioFavAddCmd.Flags().Lookup("track"))
	audioFavRemoveCmd.Flags().String("track", "", "track id")
	viper.BindPFlag("track", audioFavRemoveCmd.Flags().Lookup("track"))

	audioFavoritesCmd.AddCommand(audioFavListCmd, audioFavAddCmd, audioFavRemoveCmd)
	audioCmd.AddCommand(audioTracksCmd, audioCategoriesCmd, audioStateCmd, audioPlayCmd, audioPauseCmd, audioSeekCmd, audioVolumeCmd, audioPairCmd, audioNextCmd, audioFavoritesCmd)
}
