package client

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/charmbracelet/log"
	"github.com/krishna/relaxtech/internal/tokencache"
)

const (
	// Primary API endpoint for the Eight Sleep platform.
	defaultBaseURL = "https://client-api.8slp.net/v1"
	// OAuth authentication endpoint.
	authURL = "https://auth-api.8slp.net/v1/tokens"
)

// Client holds configuration and state for communicating with the Eight Sleep API.
type Client struct {
	Email        string
	Password     string
	UserID       string
	ClientID     string
	ClientSecret string
	DeviceID     string

	HTTP     *http.Client
	BaseURL  string
	token    string
	tokenExp time.Time
}

// New constructs a Client with the given credentials, applying defaults where needed.
func New(email, password, userID, clientID, clientSecret string) *Client {
	clientID = strings.TrimSpace(clientID)
	clientSecret = strings.TrimSpace(clientSecret)
	if clientID == "" {
		clientID = "sleep-client"
	}
	transport := &http.Transport{
		Proxy:           http.ProxyFromEnvironment,
		TLSClientConfig: &tls.Config{MinVersion: tls.VersionTLS12},
		// Force HTTP/1.1 to avoid hanging connections with the API servers
		TLSNextProto: map[string]func(string, *tls.Conn) http.RoundTripper{},
	}
	return &Client{
		Email:        email,
		Password:     password,
		UserID:       userID,
		ClientID:     clientID,
		ClientSecret: clientSecret,
		HTTP:         &http.Client{Timeout: 20 * time.Second, Transport: transport},
		BaseURL:      defaultBaseURL,
	}
}

// Authenticate obtains a bearer token, attempting the OAuth endpoint first
// and falling back to the legacy login flow if that fails.
func (c *Client) Authenticate(ctx context.Context) error {
	if err := c.authTokenEndpoint(ctx); err == nil {
		return nil
	}
	return c.authLegacyLogin(ctx)
}

// EnsureUserID fetches the current user ID from the API when it has not been set.
func (c *Client) EnsureUserID(ctx context.Context) error {
	if c.UserID != "" {
		return nil
	}
	var res struct {
		User struct {
			UserID string `json:"userId"`
		} `json:"user"`
	}
	if err := c.do(ctx, http.MethodGet, "/users/me", nil, nil, &res); err != nil {
		return err
	}
	if res.User.UserID == "" {
		return errors.New("userId not found")
	}
	c.UserID = res.User.UserID
	return nil
}

// EnsureDeviceID retrieves the device ID associated with the current user if not already known.
func (c *Client) EnsureDeviceID(ctx context.Context) (string, error) {
	if c.DeviceID != "" {
		return c.DeviceID, nil
	}
	var res struct {
		User struct {
			CurrentDevice struct {
				ID string `json:"id"`
			} `json:"currentDevice"`
		} `json:"user"`
	}
	if err := c.do(ctx, http.MethodGet, "/users/me", nil, nil, &res); err != nil {
		return "", err
	}
	if res.User.CurrentDevice.ID == "" {
		return "", errors.New("no current device id")
	}
	c.DeviceID = res.User.CurrentDevice.ID
	return c.DeviceID, nil
}

func (c *Client) authTokenEndpoint(ctx context.Context) error {
	payload := map[string]string{
		"grant_type":    "password",
		"username":      c.Email,
		"password":      c.Password,
		"client_id":     c.ClientID,
		"client_secret": c.ClientSecret,
	}
	body, _ := json.Marshal(payload)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, authURL, bytes.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.HTTP.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 300 {
		return fmt.Errorf("token auth failed: %s", resp.Status)
	}

	var res struct {
		AccessToken string `json:"access_token"`
		ExpiresIn   int    `json:"expires_in"`
		UserID      string `json:"userId"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return err
	}
	if res.AccessToken == "" {
		return errors.New("empty access token")
	}
	c.token = res.AccessToken
	if res.ExpiresIn == 0 {
		res.ExpiresIn = 3600
	}
	c.tokenExp = time.Now().Add(time.Duration(res.ExpiresIn-60) * time.Second)
	if c.UserID == "" {
		c.UserID = res.UserID
	}
	if err := tokencache.Persist(c.AuthContext(), c.token, c.tokenExp, c.UserID); err != nil {
		log.Debug("unable to persist token to cache", "error", err)
	} else {
		log.Debug("token persisted to cache", "expires_at", c.tokenExp)
	}
	return nil
}

func (c *Client) authLegacyLogin(ctx context.Context) error {
	payload := map[string]string{
		"email":    c.Email,
		"password": c.Password,
	}
	body, _ := json.Marshal(payload)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.BaseURL+"/login", bytes.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("User-Agent", "okhttp/4.9.3")
	req.Header.Set("Accept-Encoding", "gzip")
	resp, err := c.HTTP.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 300 {
		return fmt.Errorf("login failed: %s", resp.Status)
	}
	var res struct {
		Session struct {
			UserID         string `json:"userId"`
			Token          string `json:"token"`
			ExpirationDate string `json:"expirationDate"`
		} `json:"session"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return err
	}
	if res.Session.Token == "" {
		return errors.New("empty session token")
	}
	c.token = res.Session.Token
	if res.Session.ExpirationDate != "" {
		if t, err := time.Parse(time.RFC3339, res.Session.ExpirationDate); err == nil {
			c.tokenExp = t
		}
	}
	if c.tokenExp.IsZero() {
		c.tokenExp = time.Now().Add(12 * time.Hour)
	}
	if c.UserID == "" {
		c.UserID = res.Session.UserID
	}
	if err := tokencache.Persist(c.AuthContext(), c.token, c.tokenExp, c.UserID); err != nil {
		log.Debug("unable to persist token to cache", "error", err)
	} else {
		log.Debug("token persisted to cache via legacy flow", "expires_at", c.tokenExp)
	}
	return nil
}

func (c *Client) ensureToken(ctx context.Context) error {
	if c.token != "" && time.Now().Before(c.tokenExp) {
		log.Debug("reusing in-memory token", "expires_in", time.Until(c.tokenExp).Round(time.Second))
		return nil
	}
	// Attempt to restore a previously cached token. If the token turns out to be
	// expired server-side, the 401 handler in do() will clear it and re-auth.
	if cached, err := tokencache.Retrieve(c.AuthContext(), c.UserID); err == nil {
		log.Debug("restored token from cache", "expires_at", cached.ExpiresAt, "user_id", cached.UserID)
		c.token = cached.Token
		c.tokenExp = cached.ExpiresAt
		if cached.UserID != "" && c.UserID == "" {
			c.UserID = cached.UserID
		}
		return nil
	} else {
		log.Debug("cache miss for token", "reason", err)
	}
	log.Debug("performing fresh authentication")
	return c.Authenticate(ctx)
}

// requireUser guarantees that the UserID field is populated before making API calls.
func (c *Client) requireUser(ctx context.Context) error {
	if c.UserID != "" {
		return nil
	}
	return c.EnsureUserID(ctx)
}

func (c *Client) do(ctx context.Context, method, path string, query url.Values, body any, out any) error {
	return c.doWithRetry(ctx, method, path, query, body, out, false)
}

func (c *Client) doWithRetry(ctx context.Context, method, path string, query url.Values, body any, out any, retried bool) error {
	if err := c.ensureToken(ctx); err != nil {
		return err
	}
	var rdr io.Reader
	if body != nil {
		b, err := json.Marshal(body)
		if err != nil {
			return err
		}
		rdr = bytes.NewReader(b)
	}
	endpoint := c.BaseURL + path
	if len(query) > 0 {
		endpoint += "?" + query.Encode()
	}
	req, err := http.NewRequestWithContext(ctx, method, endpoint, rdr)
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+c.token)
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("User-Agent", "okhttp/4.9.3")
	req.Header.Set("Accept-Encoding", "gzip")

	resp, err := c.HTTP.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusTooManyRequests {
		return fmt.Errorf("api rate limit exceeded")
	}
	if resp.StatusCode == http.StatusUnauthorized {
		c.token = ""
		c.tokenExp = time.Time{}
		_ = tokencache.Remove(c.AuthContext())
		if err := c.Authenticate(ctx); err != nil {
			return err
		}
		if retried {
			return errors.New("authentication failed after token refresh")
		}
		return c.doWithRetry(ctx, method, path, query, body, out, true)
	}
	if resp.StatusCode >= 300 {
		return fmt.Errorf("api %s %s: %s", method, path, resp.Status)
	if out != nil {
		return json.NewDecoder(resp.Body).Decode(out)
	}
	return nil
}

// TurnOn activates the connected device.
func (c *Client) TurnOn(ctx context.Context) error {
	return c.setPower(ctx, true)
}

// TurnOff deactivates the connected device.
func (c *Client) TurnOff(ctx context.Context) error {
	return c.setPower(ctx, false)
}

func (c *Client) setPower(ctx context.Context, on bool) error {
	if err := c.requireUser(ctx); err != nil {
		return err
	}
	path := fmt.Sprintf("/users/%s/devices/power", c.UserID)
	body := map[string]bool{"on": on}
	return c.do(ctx, http.MethodPost, path, nil, body, nil)
}

func (c *Client) AuthContext() tokencache.AuthContext {
	return tokencache.AuthContext{
		BaseURL:  c.BaseURL,
		ClientID: c.ClientID,
		Email:    c.Email,
	}
}

// SetTemperature adjusts the bed temperature to the specified level (range: -100 to 100).
func (c *Client) SetTemperature(ctx context.Context, level int) error {
	if err := c.requireUser(ctx); err != nil {
		return err
	}
	if level < -100 || level > 100 {
		return fmt.Errorf("level must be between -100 and 100")
	}
	path := fmt.Sprintf("/users/%s/temperature", c.UserID)
	body := map[string]int{"currentLevel": level}
	return c.do(ctx, http.MethodPut, path, nil, body, nil)
}

// TempStatus describes the current thermal state of the bed.
type TempStatus struct {
	CurrentLevel int `json:"currentLevel"`
	CurrentState struct {
		Type string `json:"type"`
	} `json:"currentState"`
}

// GetStatus retrieves the current temperature mode and level for the user.
func (c *Client) GetStatus(ctx context.Context) (*TempStatus, error) {
	if err := c.requireUser(ctx); err != nil {
		return nil, err
	}
	path := fmt.Sprintf("/users/%s/temperature", c.UserID)
	var res TempStatus
	if err := c.do(ctx, http.MethodGet, path, nil, nil, &res); err != nil {
		return nil, err
	}
	return &res, nil
}

// SleepDay holds aggregated sleep data for a single calendar date.
type SleepDay struct {
	Date          string  `json:"day"`
	Score         float64 `json:"score"`
	Tnt           int     `json:"tnt"`
	Respiratory   float64 `json:"respiratoryRate"`
	HeartRate     float64 `json:"heartRate"`
	LatencyAsleep float64 `json:"latencyAsleepSeconds"`
	LatencyOut    float64 `json:"latencyOutSeconds"`
	Duration      float64 `json:"sleepDurationSeconds"`
	Stages        []Stage `json:"stages"`
	SleepQuality  struct {
		HRV struct {
			Score float64 `json:"score"`
		} `json:"hrv"`
		Resp struct {
			Score float64 `json:"score"`
		} `json:"respiratoryRate"`
	} `json:"sleepQualityScore"`
}

// Stage captures the duration of an individual sleep stage.
type Stage struct {
	Stage    string  `json:"stage"`
	Duration float64 `json:"duration"`
}

// GetSleepDay fetches sleep trend data for the given date string (YYYY-MM-DD) and timezone.
func (c *Client) GetSleepDay(ctx context.Context, date string, timezone string) (*SleepDay, error) {
	if err := c.requireUser(ctx); err != nil {
		return nil, err
	}
	q := url.Values{}
	q.Set("tz", timezone)
	q.Set("from", date)
	q.Set("to", date)
	q.Set("include-main", "false")
	q.Set("include-all-sessions", "true")
	q.Set("model-version", "v2")
	path := fmt.Sprintf("/users/%s/trends", c.UserID)
	var res struct {
		Days []SleepDay `json:"days"`
	}
	if err := c.do(ctx, http.MethodGet, path, q, nil, &res); err != nil {
		return nil, err
	}
	if len(res.Days) == 0 {
		return nil, fmt.Errorf("no sleep data for %s", date)
	}
	return &res.Days[0], nil
}

// AudioTrack represents metadata for a single audio content item.
type AudioTrack struct {
	ID    string `json:"id"`
	Title string `json:"title"`
	Type  string `json:"type"`
}

// ListTracks retrieves the full catalog of available audio tracks.
func (c *Client) ListTracks(ctx context.Context) ([]AudioTrack, error) {
	path := "/audio/tracks"
	var res struct {
		Tracks []AudioTrack `json:"tracks"`
	}
	if err := c.do(ctx, http.MethodGet, path, nil, nil, &res); err != nil {
		return nil, err
	}
	return res.Tracks, nil
}

// ReleaseFeature describes a feature announced in a platform release.
type ReleaseFeature struct {
	Title string `json:"title"`
	Body  string `json:"body"`
}

// ReleaseFeatures lists recently announced platform features.
func (c *Client) ReleaseFeatures(ctx context.Context) ([]ReleaseFeature, error) {
	path := "/release/features"
	var res struct {
		Features []ReleaseFeature `json:"features"`
	}
	if err := c.do(ctx, http.MethodGet, path, nil, nil, &res); err != nil {
		return nil, err
	}
	return res.Features, nil
}
