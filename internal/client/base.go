package client

import (
	"context"
	"fmt"
	"net/http"
)

// BaseActions groups endpoints for controlling the adjustable bed base.
type BaseActions struct{ c *Client }

// Base returns a handle for bed base operations.
func (c *Client) Base() *BaseActions { return &BaseActions{c: c} }

// Info retrieves current bed base status and configuration.
func (b *BaseActions) Info(ctx context.Context) (any, error) {
	if err := b.c.requireUser(ctx); err != nil {
		return nil, err
	}
	path := fmt.Sprintf("/users/%s/base", b.c.UserID)
	var res any
	err := b.c.do(ctx, http.MethodGet, path, nil, nil, &res)
	return res, err
}

// SetAngle adjusts the head and foot angles of the bed base.
func (b *BaseActions) SetAngle(ctx context.Context, head, foot int) error {
	if err := b.c.requireUser(ctx); err != nil {
		return err
	}
	path := fmt.Sprintf("/users/%s/base/angle", b.c.UserID)
	body := map[string]any{"head": head, "foot": foot}
	return b.c.do(ctx, http.MethodPost, path, nil, body, nil)
}

// Presets lists the saved position presets for the bed base.
func (b *BaseActions) Presets(ctx context.Context) (any, error) {
	if err := b.c.requireUser(ctx); err != nil {
		return nil, err
	}
	path := fmt.Sprintf("/users/%s/base/presets", b.c.UserID)
	var res any
	err := b.c.do(ctx, http.MethodGet, path, nil, nil, &res)
	return res, err
}

// RunPreset activates a named position preset on the bed base.
func (b *BaseActions) RunPreset(ctx context.Context, name string) error {
	if err := b.c.requireUser(ctx); err != nil {
		return err
	}
	path := fmt.Sprintf("/users/%s/base/presets", b.c.UserID)
	body := map[string]any{"name": name}
	return b.c.do(ctx, http.MethodPost, path, nil, body, nil)
}

// VibrationTest runs a vibration test on the device hardware.
func (b *BaseActions) VibrationTest(ctx context.Context) error {
	deviceID, err := b.c.EnsureDeviceID(ctx)
	if err != nil {
		return err
	}
	path := fmt.Sprintf("/devices/%s/vibration-test", deviceID)
	return b.c.do(ctx, http.MethodPost, path, nil, map[string]any{}, nil)
}
