// Package daemon provides a long-running scheduler that triggers device
// actions (power on/off, temperature adjustments) at configured times.
package daemon

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/krishna/relaxtech/internal/client"
)

// DefaultPIDPath is the standard location for the relaxtech daemon PID file.
const DefaultPIDPath = "~/.config/relaxtech/daemon.pid"

// TimedTask represents a single scheduled operation loaded from the
// relaxtech configuration file (~/.config/relaxtech/).
type TimedTask struct {
	Time        string `mapstructure:"time" yaml:"time"`
	Action      string `mapstructure:"action" yaml:"action"`
	Temperature string `mapstructure:"temperature" yaml:"temperature"`
}

// Scheduler is responsible for continuously checking the clock and
// dispatching TimedTask entries when their scheduled moment arrives.
type Scheduler struct {
	Tasks       []TimedTask
	DeviceAPI   *client.Client
	Location    *time.Location
	DryRun      bool
	Sync        bool
	PIDFilePath string
}

// Start launches the scheduler loop. It writes a PID file on startup and
// removes it on exit. The loop can be cancelled via the provided context or
// by sending SIGINT / SIGTERM to the process.
func (s *Scheduler) Start(ctx context.Context) error {
	if err := s.storePID(); err != nil {
		return err
	}
	defer s.clearPID()

	tick := time.NewTicker(time.Minute)
	defer tick.Stop()

	alreadyRun := make(map[string]bool)
	currentDay := time.Now().Day()

	halt := make(chan os.Signal, 1)
	signal.Notify(halt, syscall.SIGINT, syscall.SIGTERM)

	for {
		select {
		case <-ctx.Done():
			return nil
		case <-halt:
			return nil
		case ts := <-tick.C:
			// Reset the dedup map at midnight so tasks fire again.
			if ts.Day() != currentDay {
				alreadyRun = make(map[string]bool)
				currentDay = ts.Day()
			}
			if err := s.evaluate(ts, alreadyRun); err != nil {
				return err
			}
		}
	}
}

// evaluate walks through every configured task and fires the ones whose
// scheduled time falls within the current minute window.
func (s *Scheduler) evaluate(now time.Time, alreadyRun map[string]bool) error {
	for _, task := range s.Tasks {
		parsed, err := time.ParseInLocation("15:04", task.Time, s.Location)
		if err != nil {
			return fmt.Errorf("invalid time format %q: %w", task.Time, err)
		}

		target := time.Date(
			now.Year(), now.Month(), now.Day(),
			parsed.Hour(), parsed.Minute(), 0, 0, s.Location,
		)

		// Only fire if the current moment is within the one-minute window.
		if now.Before(target) || now.Sub(target) >= time.Minute {
			continue
		}

		dedupKey := target.Format("2006-01-02 15:04") + task.Action
		if alreadyRun[dedupKey] {
			continue
		}
		alreadyRun[dedupKey] = true

		if s.DryRun {
			fmt.Printf("DRY-RUN %s %s %s\n", target.Format(time.RFC3339), task.Action, task.Temperature)
			continue
		}

		if err := s.dispatch(task); err != nil {
			return err
		}
	}
	return nil
}

// dispatch sends the appropriate command to the device API based on the
// task's action field.
func (s *Scheduler) dispatch(task TimedTask) error {
	switch task.Action {
	case "on":
		return s.DeviceAPI.TurnOn(context.Background())
	case "off":
		return s.DeviceAPI.TurnOff(context.Background())
	case "temp":
		level, err := ConvertTemperature(task.Temperature)
		if err != nil {
			return err
		}
		return s.DeviceAPI.SetTemperature(context.Background(), level)
	default:
		return fmt.Errorf("unsupported action %q", task.Action)
	}
}

// storePID writes the current process ID to the configured PID file.
// It refuses to proceed if another instance appears to be running.
func (s *Scheduler) storePID() error {
	if s.PIDFilePath == "" {
		return nil
	}

	dir := filepath.Dir(s.PIDFilePath)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return fmt.Errorf("creating config directory %s: %w", dir, err)
	}

	if raw, err := os.ReadFile(s.PIDFilePath); err == nil {
		existing := strings.TrimSpace(string(raw))
		if existing != "" {
			pid, parseErr := strconv.Atoi(existing)
			if parseErr != nil || pid <= 0 {
				return fmt.Errorf("invalid daemon PID file %s", s.PIDFilePath)
			}
			process, findErr := os.FindProcess(pid)
			if findErr == nil && process.Signal(syscall.Signal(0)) == nil {
				return fmt.Errorf("relaxtech daemon already running (pid %s)", existing)
			}
			if removeErr := os.Remove(s.PIDFilePath); removeErr != nil && !errors.Is(removeErr, os.ErrNotExist) {
				return fmt.Errorf("removing stale PID file %s: %w", s.PIDFilePath, removeErr)
			}
		}
	}

	pid := fmt.Sprint(os.Getpid())
	return os.WriteFile(s.PIDFilePath, []byte(pid), 0o600)
}

// clearPID removes the PID file so that future invocations are not blocked.
func (s *Scheduler) clearPID() {
	if s.PIDFilePath != "" {
		_ = os.Remove(s.PIDFilePath)
	}
}

// ---------------------------------------------------------------------------
// Temperature conversion helpers
// ---------------------------------------------------------------------------

// ConvertTemperature interprets a user-supplied temperature string and returns
// an internal heating/cooling level in the range [-100, 100].
//
// Accepted formats:
//   - "68F"  -- Fahrenheit
//   - "20C"  -- Celsius
//   - "42"   -- raw level passed through as-is
func ConvertTemperature(input string) (int, error) {
	cleaned := strings.TrimSpace(strings.ToUpper(input))

	if strings.HasSuffix(cleaned, "F") {
		var fahrenheit float64
		if _, err := fmt.Sscanf(strings.TrimSuffix(cleaned, "F"), "%f", &fahrenheit); err != nil {
			return 0, fmt.Errorf("parsing fahrenheit value %q: %w", input, err)
		}
		return fahrenheitToLevel(fahrenheit), nil
	}

	if strings.HasSuffix(cleaned, "C") {
		var celsius float64
		if _, err := fmt.Sscanf(strings.TrimSuffix(cleaned, "C"), "%f", &celsius); err != nil {
			return 0, fmt.Errorf("parsing celsius value %q: %w", input, err)
		}
		return celsiusToLevel(celsius), nil
	}

	var rawLevel int
	if _, err := fmt.Sscanf(cleaned, "%d", &rawLevel); err == nil {
		return rawLevel, nil
	}

	return 0, fmt.Errorf("temperature %q must end with F or C, or be a numeric level", input)
}

// fahrenheitToLevel applies a linear mapping from the Fahrenheit range
// [55, 100] onto the device level range [-100, 100].
func fahrenheitToLevel(f float64) int {
	level := (f-55.0)/(100.0-55.0)*200.0 - 100.0
	return clampLevel(level)
}

// celsiusToLevel applies a linear mapping from the Celsius range
// [13, 38] (roughly 55-100 F) onto the device level range [-100, 100].
func celsiusToLevel(c float64) int {
	level := (c-13.0)/(38.0-13.0)*200.0 - 100.0
	return clampLevel(level)
}

// clampLevel restricts a floating-point level to the valid integer
// range of [-100, 100].
func clampLevel(v float64) int {
	if v < -100 {
		v = -100
	}
	if v > 100 {
		v = 100
	}
	return int(v)
}
