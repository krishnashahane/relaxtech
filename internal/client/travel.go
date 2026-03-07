package client

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
)

// TravelActions groups endpoints for managing travel plans and trip data.
type TravelActions struct{ c *Client }

// Travel returns a handle for travel-related API operations.
func (c *Client) Travel() *TravelActions { return &TravelActions{c: c} }

// Trips lists all trips associated with the user.
func (t *TravelActions) Trips(ctx context.Context) (any, error) {
	if err := t.c.requireUser(ctx); err != nil {
		return nil, err
	}
	path := fmt.Sprintf("/users/%s/travel/trips", t.c.UserID)
	var res any
	err := t.c.do(ctx, http.MethodGet, path, nil, nil, &res)
	return res, err
}

// CreateTrip registers a new trip with the provided details.
func (t *TravelActions) CreateTrip(ctx context.Context, body map[string]any) error {
	if err := t.c.requireUser(ctx); err != nil {
		return err
	}
	path := fmt.Sprintf("/users/%s/travel/trips", t.c.UserID)
	return t.c.do(ctx, http.MethodPost, path, nil, body, nil)
}

// CreatePlan adds a travel plan to the specified trip.
func (t *TravelActions) CreatePlan(ctx context.Context, tripID string, body map[string]any) error {
	if err := t.c.requireUser(ctx); err != nil {
		return err
	}
	path := fmt.Sprintf("/users/%s/travel/trips/%s/plans", t.c.UserID, tripID)
	return t.c.do(ctx, http.MethodPost, path, nil, body, nil)
}

// UpdatePlan applies modifications to an existing travel plan.
func (t *TravelActions) UpdatePlan(ctx context.Context, planID string, body map[string]any) error {
	if err := t.c.requireUser(ctx); err != nil {
		return err
	}
	path := fmt.Sprintf("/users/%s/travel/plans/%s", t.c.UserID, planID)
	return t.c.do(ctx, http.MethodPatch, path, nil, body, nil)
}

// DeleteTrip removes a trip and its associated plans.
func (t *TravelActions) DeleteTrip(ctx context.Context, tripID string) error {
	if err := t.c.requireUser(ctx); err != nil {
		return err
	}
	path := fmt.Sprintf("/users/%s/travel/trips/%s", t.c.UserID, tripID)
	return t.c.do(ctx, http.MethodDelete, path, nil, nil, nil)
}

// Plans lists all plans within the specified trip.
func (t *TravelActions) Plans(ctx context.Context, tripID string) (any, error) {
	if err := t.c.requireUser(ctx); err != nil {
		return nil, err
	}
	path := fmt.Sprintf("/users/%s/travel/trips/%s/plans", t.c.UserID, tripID)
	var res any
	err := t.c.do(ctx, http.MethodGet, path, nil, nil, &res)
	return res, err
}

// PlanTasks retrieves the task list for a specific travel plan.
func (t *TravelActions) PlanTasks(ctx context.Context, planID string) (any, error) {
	if err := t.c.requireUser(ctx); err != nil {
		return nil, err
	}
	path := fmt.Sprintf("/users/%s/travel/plans/%s/tasks", t.c.UserID, planID)
	var res any
	err := t.c.do(ctx, http.MethodGet, path, nil, nil, &res)
	return res, err
}

// AirportSearch finds airports matching the given search query.
func (t *TravelActions) AirportSearch(ctx context.Context, query string) (any, error) {
	q := url.Values{"query": []string{query}}
	var res any
	err := t.c.do(ctx, http.MethodGet, "/travel/airport-search", q, nil, &res)
	return res, err
}

// FlightStatus retrieves status information for the given flight number.
func (t *TravelActions) FlightStatus(ctx context.Context, flight string) (any, error) {
	q := url.Values{"flightNumber": []string{flight}}
	var res any
	err := t.c.do(ctx, http.MethodGet, "/travel/flight-status", q, nil, &res)
	return res, err
}
