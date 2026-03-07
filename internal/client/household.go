package client

import (
	"context"
	"fmt"
	"net/http"
)

// HouseholdActions groups endpoints for managing the household account.
type HouseholdActions struct{ c *Client }

// Household returns a handle for household-related API operations.
func (c *Client) Household() *HouseholdActions { return &HouseholdActions{c: c} }

// Summary retrieves a high-level overview of the household.
func (h *HouseholdActions) Summary(ctx context.Context) (any, error) {
	if err := h.c.requireUser(ctx); err != nil {
		return nil, err
	}
	path := fmt.Sprintf("/household/users/%s/summary", h.c.UserID)
	var res any
	err := h.c.do(ctx, http.MethodGet, path, nil, nil, &res)
	return res, err
}

// Schedule returns the household sleep schedule.
func (h *HouseholdActions) Schedule(ctx context.Context) (any, error) {
	if err := h.c.requireUser(ctx); err != nil {
		return nil, err
	}
	path := fmt.Sprintf("/household/users/%s/schedule", h.c.UserID)
	var res any
	err := h.c.do(ctx, http.MethodGet, path, nil, nil, &res)
	return res, err
}

// CurrentSet fetches the active configuration set for the household.
func (h *HouseholdActions) CurrentSet(ctx context.Context) (any, error) {
	if err := h.c.requireUser(ctx); err != nil {
		return nil, err
	}
	path := fmt.Sprintf("/household/users/%s/current-set", h.c.UserID)
	var res any
	err := h.c.do(ctx, http.MethodGet, path, nil, nil, &res)
	return res, err
}

// Invitations lists pending household invitations.
func (h *HouseholdActions) Invitations(ctx context.Context) (any, error) {
	if err := h.c.requireUser(ctx); err != nil {
		return nil, err
	}
	path := fmt.Sprintf("/household/users/%s/invitations", h.c.UserID)
	var res any
	err := h.c.do(ctx, http.MethodGet, path, nil, nil, &res)
	return res, err
}

// Devices lists all devices registered to the household.
func (h *HouseholdActions) Devices(ctx context.Context) (any, error) {
	if err := h.c.requireUser(ctx); err != nil {
		return nil, err
	}
	path := fmt.Sprintf("/household/users/%s/devices", h.c.UserID)
	var res any
	err := h.c.do(ctx, http.MethodGet, path, nil, nil, &res)
	return res, err
}

// Users lists the members of the household.
func (h *HouseholdActions) Users(ctx context.Context) (any, error) {
	if err := h.c.requireUser(ctx); err != nil {
		return nil, err
	}
	path := fmt.Sprintf("/household/users/%s/users", h.c.UserID)
	var res any
	err := h.c.do(ctx, http.MethodGet, path, nil, nil, &res)
	return res, err
}

// Guests lists guest accounts associated with the household.
func (h *HouseholdActions) Guests(ctx context.Context) (any, error) {
	if err := h.c.requireUser(ctx); err != nil {
		return nil, err
	}
	path := fmt.Sprintf("/household/users/%s/guests", h.c.UserID)
	var res any
	err := h.c.do(ctx, http.MethodGet, path, nil, nil, &res)
	return res, err
}
