package cmd

import (
	"path/filepath"
	"testing"
	"time"

	"github.com/99designs/keyring"
	"github.com/spf13/viper"

	"github.com/krishna/relaxtech/internal/client"
	"github.com/krishna/relaxtech/internal/tokencache"
)

// setupTemporaryKeyring configures an ephemeral file-backed keyring for testing.
func setupTemporaryKeyring(t *testing.T) func() {
	t.Helper()
	dir := t.TempDir()
	cleanup := tokencache.OverrideRingOpenerForTest(func() (keyring.Keyring, error) {
		return keyring.Open(keyring.Config{
			ServiceName:      "relaxtech-test",
			AllowedBackends:  []keyring.BackendType{keyring.FileBackend},
			FileDir:          filepath.Join(dir, "keyring"),
			FilePasswordFunc: func(_ string) (string, error) { return "test-pass", nil },
		})
	})
	t.Cleanup(cleanup)
	return cleanup
}

// clearViper resets all viper state between test runs.
func clearViper(t *testing.T) {
	t.Helper()
	viper.Reset()
}

func TestEnsureCredentialsSucceedsWithStoredToken(t *testing.T) {
	setupTemporaryKeyring(t)
	clearViper(t)

	// Persist a token without providing credentials directly.
	cl := client.New("", "", "", "", "")
	if err := tokencache.Persist(cl.AuthContext(), "tok", time.Now().Add(time.Hour), "stored-user"); err != nil {
		t.Fatalf("persist token: %v", err)
	}

	if err := ensureCredentials(); err != nil {
		t.Fatalf("ensureCredentials should succeed with stored token: %v", err)
	}
	if got := viper.GetString("user_id"); got != "stored-user" {
		t.Fatalf("user_id not propagated from stored token, got %q", got)
	}
}

func TestEnsureCredentialsFailsWithoutTokenOrCredentials(t *testing.T) {
	setupTemporaryKeyring(t)
	clearViper(t)

	err := ensureCredentials()
	if err == nil {
		t.Fatalf("expected error for missing credentials")
	}
}
