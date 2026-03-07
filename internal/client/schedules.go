package client

import (
	"context"
	"fmt"
	"net/http"
)

// TemperatureSchedule defines a recurring temperature adjustment rule.
type TemperatureSchedule struct {
	ID         string `json:"id"`
	StartTime  string `json:"startTime"`
	Level      int    `json:"level"`
	DaysOfWeek []int  `json:"daysOfWeek"`
	Enabled    bool   `json:"enabled"`
}

// ListSchedules retrieves all temperature schedules for the current user.
func (c *Client) ListSchedules(ctx context.Context) ([]TemperatureSchedule, error) {
	if err := c.requireUser(ctx); err != nil {
		return nil, err
	}
	path := fmt.Sprintf("/users/%s/temperature/schedules", c.UserID)
	var res struct {
		Schedules []TemperatureSchedule `json:"schedules"`
	}
	if err := c.do(ctx, http.MethodGet, path, nil, nil, &res); err != nil {
		return nil, err
	}
	return res.Schedules, nil
}

// CreateSchedule adds a new temperature schedule and returns the created resource.
func (c *Client) CreateSchedule(ctx context.Context, s TemperatureSchedule) (*TemperatureSchedule, error) {
	if err := c.requireUser(ctx); err != nil {
		return nil, err
	}
	path := fmt.Sprintf("/users/%s/temperature/schedules", c.UserID)
	var res struct {
		Schedule TemperatureSchedule `json:"schedule"`
	}
	if err := c.do(ctx, http.MethodPost, path, nil, s, &res); err != nil {
		return nil, err
	}
	return &res.Schedule, nil
}

// UpdateSchedule applies partial modifications to an existing schedule by ID.
func (c *Client) UpdateSchedule(ctx context.Context, id string, patch map[string]any) (*TemperatureSchedule, error) {
	if err := c.requireUser(ctx); err != nil {
		return nil, err
	}
	path := fmt.Sprintf("/users/%s/temperature/schedules/%s", c.UserID, id)
	var res struct {
		Schedule TemperatureSchedule `json:"schedule"`
	}
	if err := c.do(ctx, http.MethodPatch, path, nil, patch, &res); err != nil {
		return nil, err
	}
	return &res.Schedule, nil
}

// DeleteSchedule removes a temperature schedule by its ID.
func (c *Client) DeleteSchedule(ctx context.Context, id string) error {
	if err := c.requireUser(ctx); err != nil {
		return err
	}
	path := fmt.Sprintf("/users/%s/temperature/schedules/%s", c.UserID, id)
	return c.do(ctx, http.MethodDelete, path, nil, nil, nil)
}
