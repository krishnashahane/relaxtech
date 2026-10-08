package tokencache

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/99designs/keyring"
)

// setupTestRing configures a temporary file-backed keyring for the
// duration of the test.
func setupTestRing(t *testing.T) {
	t.Helper()
	dir := t.TempDir()
	previous := ringOpener
	ringOpener = func() (keyring.Keyring, error) {
		return keyring.Open(keyring.Config{
			ServiceName:      appName + "-test",
			AllowedBackends:  []keyring.BackendType{keyring.FileBackend},
			FileDir:          filepath.Join(dir, "keyring"),
			FilePasswordFunc: func(_ string) (string, error) { return "test-pass", nil },
		})
	}
	t.Cleanup(func() { ringOpener = previous })
}

func TestPersistAndRetrieveRoundTrip(t *testing.T) {
	setupTestRing(t)

	ctx := AuthContext{BaseURL: "https://api.example.com", ClientID: "client-1", Email: "User@Example.com"}
	expiry := time.Now().Add(time.Hour)

	if err := Persist(ctx, "token-123", expiry, "user-1"); err != nil {
		t.Fatalf("Persist failed: %v", err)
	}

	result, err := Retrieve(ctx, "user-1")
	if err != nil {
		t.Fatalf("Retrieve failed: %v", err)
	}
	if result.Token != "token-123" {
		t.Errorf("got token %q, expected token-123", result.Token)
	}
	if !result.ExpiresAt.Equal(expiry) {
		t.Errorf("got expiry %v, expected %v", result.ExpiresAt, expiry)
	}
	if result.UserID != "user-1" {
		t.Errorf("got user ID %q, expected user-1", result.UserID)
	}
}

func TestRetrieveRejectsMismatchedUser(t *testing.T) {
	setupTestRing(t)
	ctx := AuthContext{BaseURL: "https://api.example.com", ClientID: "client-1"}
	if err := Persist(ctx, "token", time.Now().Add(time.Hour), "user-a"); err != nil {
		t.Fatalf("Persist failed: %v", err)
	}
	if _, err := Retrieve(ctx, "user-b"); err != keyring.ErrKeyNotFound {
		t.Fatalf("expected ErrKeyNotFound for wrong user, got %v", err)
	}
}

func TestRetrieveEvictsExpiredToken(t *testing.T) {
	setupTestRing(t)
	ctx := AuthContext{BaseURL: "https://api.example.com", ClientID: "client-1"}
	if err := Persist(ctx, "stale", time.Now().Add(-time.Minute), "user-1"); err != nil {
		t.Fatalf("Persist failed: %v", err)
	}
	if _, err := Retrieve(ctx, "user-1"); err != keyring.ErrKeyNotFound {
		t.Fatalf("expected ErrKeyNotFound for expired token, got %v", err)
	}
	// Confirm the entry was actually cleaned up.
	if _, err := Retrieve(ctx, "user-1"); err != keyring.ErrKeyNotFound {
		t.Fatalf("expected ErrKeyNotFound after eviction, got %v", err)
	}
}

func TestRemoveIgnoresAbsentEntry(t *testing.T) {
	setupTestRing(t)
	ctx := AuthContext{BaseURL: "https://api.example.com", ClientID: "client-1"}
	if err := Remove(ctx); err != nil {
		t.Fatalf("Remove on absent entry failed: %v", err)
	}
}

func TestTokenIsolationByAuthContext(t *testing.T) {
	setupTestRing(t)
	ctxA := AuthContext{BaseURL: "https://api.example.com", ClientID: "client-1", Email: "a@example.com"}
	ctxB := AuthContext{BaseURL: "https://api.example.com", ClientID: "client-2", Email: "a@example.com"}
	ctxC := AuthContext{BaseURL: "https://api.example.com", ClientID: "client-1", Email: "b@example.com"}
	ctxD := AuthContext{BaseURL: "https://api.example.com", ClientID: "client-1", Email: ""}

	if err := Persist(ctxA, "token-a", time.Now().Add(time.Hour), "user-a"); err != nil {
		t.Fatalf("Persist A: %v", err)
	}
	if err := Persist(ctxB, "token-b", time.Now().Add(time.Hour), "user-b"); err != nil {
		t.Fatalf("Persist B: %v", err)
	}
	if err := Persist(ctxC, "token-c", time.Now().Add(time.Hour), "user-c"); err != nil {
		t.Fatalf("Persist C: %v", err)
	}
	if err := Persist(ctxD, "token-d", time.Now().Add(time.Hour), "user-d"); err != nil {
		t.Fatalf("Persist D: %v", err)
	}

	if result, _ := Retrieve(ctxA, "user-a"); result.Token != "token-a" {
		t.Errorf("Retrieve A: got %q, want token-a", result.Token)
	}
	if _, err := Retrieve(ctxA, "user-b"); err != keyring.ErrKeyNotFound {
		t.Errorf("Retrieve A with wrong user should fail, got %v", err)
	}
	if result, _ := Retrieve(ctxB, "user-b"); result.Token != "token-b" {
		t.Errorf("Retrieve B: got %q, want token-b", result.Token)
	}
	if result, _ := Retrieve(ctxC, "user-c"); result.Token != "token-c" {
		t.Errorf("Retrieve C: got %q, want token-c", result.Token)
	}
	if result, _ := Retrieve(ctxD, ""); result.Token != "token-d" {
		t.Errorf("Retrieve D: got %q, want token-d", result.Token)
	}
}

func TestRemoveOnlyAffectsTargetContext(t *testing.T) {
	setupTestRing(t)
	ctxA := AuthContext{BaseURL: "https://api.example.com", ClientID: "client-1", Email: "a@example.com"}
	ctxB := AuthContext{BaseURL: "https://api.example.com", ClientID: "client-2", Email: "a@example.com"}

	if err := Persist(ctxA, "token-a", time.Now().Add(time.Hour), "user-a"); err != nil {
		t.Fatalf("Persist A: %v", err)
	}
	if err := Persist(ctxB, "token-b", time.Now().Add(time.Hour), "user-b"); err != nil {
		t.Fatalf("Persist B: %v", err)
	}

	if err := Remove(ctxA); err != nil {
		t.Fatalf("Remove A: %v", err)
	}
	if _, err := Retrieve(ctxA, "user-a"); err != keyring.ErrKeyNotFound {
		t.Fatalf("expected A to be removed, got %v", err)
	}
	if result, err := Retrieve(ctxB, "user-b"); err != nil || result.Token != "token-b" {
		t.Fatalf("B should be unaffected, got %v err %v", result, err)
	}
}

func TestBuildKeyNormalization(t *testing.T) {
	k1 := buildKey(AuthContext{BaseURL: "https://API.example.com/", ClientID: "id", Email: "User@Example.com "})
	k2 := buildKey(AuthContext{BaseURL: "https://api.example.com", ClientID: "id", Email: "user@example.com"})
	if k1 != k2 {
		t.Fatalf("keys should be equal after normalization; got %q vs %q", k1, k2)
	}
}

func TestBuildKeyNormalizesBlankEmail(t *testing.T) {
	k1 := buildKey(AuthContext{BaseURL: "https://api.example.com", ClientID: "id", Email: ""})
	k2 := buildKey(AuthContext{BaseURL: "https://api.example.com/", ClientID: "id", Email: " "})
	if k1 != k2 {
		t.Fatalf("blank emails should normalize identically; got %q vs %q", k1, k2)
	}
}

func TestRetrieveWithoutEmailFindsUniqueMatch(t *testing.T) {
	setupTestRing(t)
	ctx := AuthContext{BaseURL: "https://api.example.com", ClientID: "client-1", Email: "user@example.com"}
	if err := Persist(ctx, "tok", time.Now().Add(time.Hour), "user-1"); err != nil {
		t.Fatalf("Persist: %v", err)
	}

	// Omit email -- should still locate the sole matching token.
	noEmail := AuthContext{BaseURL: ctx.BaseURL, ClientID: ctx.ClientID}
	result, err := Retrieve(noEmail, "user-1")
	if err != nil {
		t.Fatalf("Retrieve without email: %v", err)
	}
	if result.Token != "tok" {
		t.Fatalf("token mismatch: got %q", result.Token)
	}
}

func TestRetrieveWithoutEmailFailsOnAmbiguity(t *testing.T) {
	setupTestRing(t)
	base := AuthContext{BaseURL: "https://api.example.com", ClientID: "client-1"}
	if err := Persist(AuthContext{BaseURL: base.BaseURL, ClientID: base.ClientID, Email: "a@example.com"}, "ta", time.Now().Add(time.Hour), "ua"); err != nil {
		t.Fatalf("persist a: %v", err)
	}
	if err := Persist(AuthContext{BaseURL: base.BaseURL, ClientID: base.ClientID, Email: "b@example.com"}, "tb", time.Now().Add(time.Hour), "ub"); err != nil {
		t.Fatalf("persist b: %v", err)
	}
	if _, err := Retrieve(base, ""); err != keyring.ErrKeyNotFound {
		t.Fatalf("expected not found with ambiguous matches, got %v", err)
	}
}

func TestFallbackPasswordFunc(t *testing.T) {
	t.Setenv("RELAXTECH_KEYRING_PASSWORD", "strong-test-password")
	pw, err := fallbackPassword("ignored")
	if err != nil {
		t.Fatalf("fallbackPassword returned error: %v", err)
	}
	if pw != "strong-test-password" {
		t.Fatalf("password = %q, want configured password", pw)
	}

	t.Setenv("RELAXTECH_KEYRING_PASSWORD", "")
	if _, err := fallbackPassword("ignored"); err == nil {
		t.Fatal("expected error when RELAXTECH_KEYRING_PASSWORD is unset")
	}
}
