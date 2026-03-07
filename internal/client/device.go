package client

import (
	"context"
	"fmt"
	"net/http"
)

// DeviceActions groups endpoints for querying hardware device information.
type DeviceActions struct{ c *Client }

// Device returns a handle for device-related API operations.
func (c *Client) Device() *DeviceActions { return &DeviceActions{c: c} }

// Info retrieves general information about the connected device.
func (d *DeviceActions) Info(ctx context.Context) (any, error) {
	deviceID, err := d.c.EnsureDeviceID(ctx)
	if err != nil {
		return nil, err
	}
	path := fmt.Sprintf("/devices/%s", deviceID)
	var res any
	err = d.c.do(ctx, http.MethodGet, path, nil, nil, &res)
	return res, err
}

// Peripherals lists the peripheral components attached to the device.
func (d *DeviceActions) Peripherals(ctx context.Context) (any, error) {
	deviceID, err := d.c.EnsureDeviceID(ctx)
	if err != nil {
		return nil, err
	}
	path := fmt.Sprintf("/devices/%s/peripherals", deviceID)
	var res any
	err = d.c.do(ctx, http.MethodGet, path, nil, nil, &res)
	return res, err
}

// Owner returns ownership details for the device.
func (d *DeviceActions) Owner(ctx context.Context) (any, error) {
	deviceID, err := d.c.EnsureDeviceID(ctx)
	if err != nil {
		return nil, err
	}
	path := fmt.Sprintf("/devices/%s/owner", deviceID)
	var res any
	err = d.c.do(ctx, http.MethodGet, path, nil, nil, &res)
	return res, err
}

// Warranty fetches warranty information for the device.
func (d *DeviceActions) Warranty(ctx context.Context) (any, error) {
	deviceID, err := d.c.EnsureDeviceID(ctx)
	if err != nil {
		return nil, err
	}
	path := fmt.Sprintf("/devices/%s/warranty", deviceID)
	var res any
	err = d.c.do(ctx, http.MethodGet, path, nil, nil, &res)
	return res, err
}

// Online checks whether the device is currently connected and reachable.
func (d *DeviceActions) Online(ctx context.Context) (any, error) {
	deviceID, err := d.c.EnsureDeviceID(ctx)
	if err != nil {
		return nil, err
	}
	path := fmt.Sprintf("/devices/%s/online", deviceID)
	var res any
	err = d.c.do(ctx, http.MethodGet, path, nil, nil, &res)
	return res, err
}

// PrimingTasks retrieves the list of priming tasks for the device.
func (d *DeviceActions) PrimingTasks(ctx context.Context) (any, error) {
	deviceID, err := d.c.EnsureDeviceID(ctx)
	if err != nil {
		return nil, err
	}
	path := fmt.Sprintf("/devices/%s/priming/tasks", deviceID)
	var res any
	err = d.c.do(ctx, http.MethodGet, path, nil, nil, &res)
	return res, err
}

// PrimingSchedule fetches the scheduled priming configuration for the device.
func (d *DeviceActions) PrimingSchedule(ctx context.Context) (any, error) {
	deviceID, err := d.c.EnsureDeviceID(ctx)
	if err != nil {
		return nil, err
	}
	path := fmt.Sprintf("/devices/%s/priming/schedule", deviceID)
	var res any
	err = d.c.do(ctx, http.MethodGet, path, nil, nil, &res)
	return res, err
}
