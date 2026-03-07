package client

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
)

// TempModes provides access to specialized temperature mode operations such as nap and hot flash.
type TempModes struct{ c *Client }

// TempModes returns a handle for temperature mode endpoints.
func (c *Client) TempModes() *TempModes { return &TempModes{c: c} }

// NapActivate starts nap mode for the user.
func (t *TempModes) NapActivate(ctx context.Context) error {
	return t.simplePost(ctx, "/temperature/nap-mode/activate")
}

// NapDeactivate stops nap mode for the user.
func (t *TempModes) NapDeactivate(ctx context.Context) error {
	return t.simplePost(ctx, "/temperature/nap-mode/deactivate")
}

// NapExtend prolongs the current nap session.
func (t *TempModes) NapExtend(ctx context.Context) error {
	return t.simplePost(ctx, "/temperature/nap-mode/extend")
}

// NapStatus retrieves the current state of nap mode.
func (t *TempModes) NapStatus(ctx context.Context, out any) error {
	return t.simpleGet(ctx, "/temperature/nap-mode/status", out)
}

// HotFlashActivate enables the hot flash mitigation mode.
func (t *TempModes) HotFlashActivate(ctx context.Context) error {
	return t.simplePost(ctx, "/temperature/hot-flash-mode/activate")
}

// HotFlashDeactivate disables the hot flash mitigation mode.
func (t *TempModes) HotFlashDeactivate(ctx context.Context) error {
	return t.simplePost(ctx, "/temperature/hot-flash-mode/deactivate")
}

// HotFlashStatus returns the current hot flash mode configuration.
func (t *TempModes) HotFlashStatus(ctx context.Context, out any) error {
	return t.simpleGet(ctx, "/temperature/hot-flash-mode", out)
}

// TempEvents retrieves temperature event history within an optional date range.
func (t *TempModes) TempEvents(ctx context.Context, from, to string, out any) error {
	if err := t.c.requireUser(ctx); err != nil {
		return err
	}
	q := url.Values{}
	if from != "" {
		q.Set("from", from)
	}
	if to != "" {
		q.Set("to", to)
	}
	path := fmt.Sprintf("/users/%s/temp-events", t.c.UserID)
	return t.c.do(ctx, http.MethodGet, path, q, nil, out)
}

func (t *TempModes) simplePost(ctx context.Context, suffix string) error {
	if err := t.c.requireUser(ctx); err != nil {
		return err
	}
	path := fmt.Sprintf("/users/%s%s", t.c.UserID, suffix)
	return t.c.do(ctx, http.MethodPost, path, nil, map[string]string{}, nil)
}

func (t *TempModes) simpleGet(ctx context.Context, suffix string, out any) error {
	if err := t.c.requireUser(ctx); err != nil {
		return err
	}
	path := fmt.Sprintf("/users/%s%s", t.c.UserID, suffix)
	return t.c.do(ctx, http.MethodGet, path, nil, nil, out)
}
