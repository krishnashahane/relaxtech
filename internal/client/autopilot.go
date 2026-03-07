package client

import (
	"context"
	"fmt"
	"net/http"
)

// AutopilotActions provides methods for managing the autopilot feature.
type AutopilotActions struct{ c *Client }

// Autopilot returns a handle for autopilot-related API operations.
func (c *Client) Autopilot() *AutopilotActions { return &AutopilotActions{c: c} }

// Details retrieves the current autopilot configuration for the user.
func (a *AutopilotActions) Details(ctx context.Context) (any, error) {
	if err := a.c.requireUser(ctx); err != nil {
		return nil, err
	}
	path := fmt.Sprintf("/users/%s/autopilotDetails", a.c.UserID)
	var res any
	err := a.c.do(ctx, http.MethodGet, path, nil, nil, &res)
	return res, err
}

// History returns historical autopilot adjustment data.
func (a *AutopilotActions) History(ctx context.Context) (any, error) {
	if err := a.c.requireUser(ctx); err != nil {
		return nil, err
	}
	path := fmt.Sprintf("/users/%s/autopilot-history", a.c.UserID)
	var res any
	err := a.c.do(ctx, http.MethodGet, path, nil, nil, &res)
	return res, err
}

// Recap fetches the autopilot recap summary.
func (a *AutopilotActions) Recap(ctx context.Context) (any, error) {
	if err := a.c.requireUser(ctx); err != nil {
		return nil, err
	}
	path := fmt.Sprintf("/users/%s/autopilotDetails/autopilotRecap", a.c.UserID)
	var res any
	err := a.c.do(ctx, http.MethodGet, path, nil, nil, &res)
	return res, err
}

// SetLevelSuggestions toggles the automatic temperature level suggestion feature.
func (a *AutopilotActions) SetLevelSuggestions(ctx context.Context, enabled bool) error {
	if err := a.c.requireUser(ctx); err != nil {
		return err
	}
	path := fmt.Sprintf("/users/%s/level-suggestions-mode", a.c.UserID)
	body := map[string]any{"enabled": enabled}
	return a.c.do(ctx, http.MethodPost, path, nil, body, nil)
}

// SetSnoreMitigation toggles the snore mitigation feature within autopilot.
func (a *AutopilotActions) SetSnoreMitigation(ctx context.Context, enabled bool) error {
	if err := a.c.requireUser(ctx); err != nil {
		return err
	}
	path := fmt.Sprintf("/users/%s/autopilotDetails/snoringMitigation", a.c.UserID)
	body := map[string]any{"enabled": enabled}
	return a.c.do(ctx, http.MethodPost, path, nil, body, nil)
}
