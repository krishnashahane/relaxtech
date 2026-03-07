package client

import (
	"context"
	"fmt"
	"net/http"
)

// Presence holds the bed occupancy status for a user.
type Presence struct {
	Present bool `json:"presence"`
}

// GetPresence checks whether the user is currently detected in bed.
func (c *Client) GetPresence(ctx context.Context) (bool, error) {
	if err := c.requireUser(ctx); err != nil {
		return false, err
	}
	path := fmt.Sprintf("/users/%s/presence", c.UserID)
	var res Presence
	if err := c.do(ctx, http.MethodGet, path, nil, nil, &res); err != nil {
		return false, err
	}
	return res.Present, nil
}
