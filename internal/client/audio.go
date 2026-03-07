package client

import (
	"context"
	"fmt"
	"net/http"
)

// AudioActions groups endpoints for audio playback and track management.
type AudioActions struct{ c *Client }

// Audio returns a handle for audio-related API operations.
func (c *Client) Audio() *AudioActions { return &AudioActions{c: c} }

// Tracks lists audio tracks available to the user.
func (a *AudioActions) Tracks(ctx context.Context) ([]AudioTrack, error) {
	if err := a.c.requireUser(ctx); err != nil {
		return nil, err
	}
	path := fmt.Sprintf("/users/%s/audio/tracks", a.c.UserID)
	var res struct {
		Tracks []AudioTrack `json:"tracks"`
	}
	if err := a.c.do(ctx, http.MethodGet, path, nil, nil, &res); err != nil {
		return nil, err
	}
	return res.Tracks, nil
}

// Categories retrieves the available audio content categories.
func (a *AudioActions) Categories(ctx context.Context) (any, error) {
	path := "/audio/categories"
	var res any
	err := a.c.do(ctx, http.MethodGet, path, nil, nil, &res)
	return res, err
}

// PlayerState returns the current state of the audio player.
func (a *AudioActions) PlayerState(ctx context.Context) (any, error) {
	if err := a.c.requireUser(ctx); err != nil {
		return nil, err
	}
	path := fmt.Sprintf("/users/%s/audio/player/state", a.c.UserID)
	var res any
	err := a.c.do(ctx, http.MethodGet, path, nil, nil, &res)
	return res, err
}

// Play starts playback, optionally for a specific track when trackID is non-empty.
func (a *AudioActions) Play(ctx context.Context, trackID string) error {
	if err := a.c.requireUser(ctx); err != nil {
		return err
	}
	path := fmt.Sprintf("/users/%s/audio/player", a.c.UserID)
	body := map[string]any{"action": "play"}
	if trackID != "" {
		body["trackId"] = trackID
	}
	return a.c.do(ctx, http.MethodPost, path, nil, body, nil)
}

// Pause halts audio playback.
func (a *AudioActions) Pause(ctx context.Context) error {
	if err := a.c.requireUser(ctx); err != nil {
		return err
	}
	path := fmt.Sprintf("/users/%s/audio/player", a.c.UserID)
	body := map[string]any{"action": "pause"}
	return a.c.do(ctx, http.MethodPost, path, nil, body, nil)
}

// Seek jumps to the given position in the current track (in milliseconds).
func (a *AudioActions) Seek(ctx context.Context, positionMs int) error {
	if err := a.c.requireUser(ctx); err != nil {
		return err
	}
	path := fmt.Sprintf("/users/%s/audio/player/seek", a.c.UserID)
	body := map[string]any{"position": positionMs}
	return a.c.do(ctx, http.MethodPost, path, nil, body, nil)
}

// Volume sets the audio playback volume to the specified level.
func (a *AudioActions) Volume(ctx context.Context, level int) error {
	if err := a.c.requireUser(ctx); err != nil {
		return err
	}
	path := fmt.Sprintf("/users/%s/audio/player/volume", a.c.UserID)
	body := map[string]any{"level": level}
	return a.c.do(ctx, http.MethodPost, path, nil, body, nil)
}

// Pair initiates Bluetooth pairing for the audio subsystem on the device.
func (a *AudioActions) Pair(ctx context.Context) error {
	deviceID, err := a.c.EnsureDeviceID(ctx)
	if err != nil {
		return err
	}
	path := fmt.Sprintf("/devices/%s/audio/player/pair", deviceID)
	return a.c.do(ctx, http.MethodPost, path, nil, map[string]any{}, nil)
}

// RecommendedNext returns the next recommended track for the user.
func (a *AudioActions) RecommendedNext(ctx context.Context) (any, error) {
	if err := a.c.requireUser(ctx); err != nil {
		return nil, err
	}
	path := fmt.Sprintf("/users/%s/audio/tracks/recommended-next-track", a.c.UserID)
	var res any
	err := a.c.do(ctx, http.MethodGet, path, nil, nil, &res)
	return res, err
}

// Favorites lists the user's favorited audio tracks.
func (a *AudioActions) Favorites(ctx context.Context) (any, error) {
	if err := a.c.requireUser(ctx); err != nil {
		return nil, err
	}
	path := fmt.Sprintf("/users/%s/audio/tracks/favorites", a.c.UserID)
	var res any
	err := a.c.do(ctx, http.MethodGet, path, nil, nil, &res)
	return res, err
}

// AddFavorite marks a track as a favorite for the user.
func (a *AudioActions) AddFavorite(ctx context.Context, trackID string) error {
	if err := a.c.requireUser(ctx); err != nil {
		return err
	}
	path := fmt.Sprintf("/users/%s/audio/tracks/favorites", a.c.UserID)
	body := map[string]any{"trackId": trackID}
	return a.c.do(ctx, http.MethodPost, path, nil, body, nil)
}

// RemoveFavorite removes a track from the user's favorites.
func (a *AudioActions) RemoveFavorite(ctx context.Context, trackID string) error {
	if err := a.c.requireUser(ctx); err != nil {
		return err
	}
	path := fmt.Sprintf("/users/%s/audio/tracks/favorites/%s", a.c.UserID, trackID)
	return a.c.do(ctx, http.MethodDelete, path, nil, nil, nil)
}
