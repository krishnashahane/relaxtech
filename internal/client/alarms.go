package client

import (
	"context"
	"fmt"
	"net/http"
)

// AlarmActions provides operations for interacting with active alarms.
type AlarmActions struct {
	c *Client
}

// Alarms returns a handle for alarm interaction endpoints.
func (c *Client) Alarms() *AlarmActions { return &AlarmActions{c: c} }

// Snooze pauses the specified alarm temporarily.
func (a *AlarmActions) Snooze(ctx context.Context, alarmID string) error {
	if err := a.c.requireUser(ctx); err != nil {
		return err
	}
	path := fmt.Sprintf("/users/%s/alarms/%s/snooze", a.c.UserID, alarmID)
	return a.c.do(ctx, http.MethodPost, path, nil, map[string]string{}, nil)
}

// Dismiss stops the specified alarm from ringing.
func (a *AlarmActions) Dismiss(ctx context.Context, alarmID string) error {
	if err := a.c.requireUser(ctx); err != nil {
		return err
	}
	path := fmt.Sprintf("/users/%s/alarms/%s/dismiss", a.c.UserID, alarmID)
	return a.c.do(ctx, http.MethodPost, path, nil, map[string]string{}, nil)
}

// DismissAll silences every currently active alarm.
func (a *AlarmActions) DismissAll(ctx context.Context) error {
	if err := a.c.requireUser(ctx); err != nil {
		return err
	}
	path := fmt.Sprintf("/users/%s/alarms/active/dismiss-all", a.c.UserID)
	return a.c.do(ctx, http.MethodPost, path, nil, map[string]string{}, nil)
}

// VibrationTest triggers a test vibration on the device.
func (a *AlarmActions) VibrationTest(ctx context.Context) error {
	if err := a.c.requireUser(ctx); err != nil {
		return err
	}
	path := fmt.Sprintf("/users/%s/vibration-test", a.c.UserID)
	return a.c.do(ctx, http.MethodPost, path, nil, map[string]string{}, nil)
}
