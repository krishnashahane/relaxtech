package tokencache

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/99designs/keyring"
	"github.com/charmbracelet/log"
)

const (
	appName       = "relaxtech"
	storedTokenID = "oauth-token"
)

// StoredToken holds a serialized authentication token along with its
// expiration and associated user identifier.
type StoredToken struct {
	Token     string    `json:"token"`
	ExpiresAt time.Time `json:"expires_at"`
	UserID    string    `json:"user_id,omitempty"`
}

// AuthContext identifies which account and environment a cached token
// belongs to. Tokens are partitioned by API endpoint, client ID, and
// email address so that multiple accounts can coexist.
type AuthContext struct {
	BaseURL  string
	ClientID string
	Email    string
}

var ringOpener = productionRingOpener

// OverrideRingOpenerForTest replaces the keyring opener function for
// testing purposes. Returns a function that restores the original.
// Must not be used concurrently.
func OverrideRingOpenerForTest(opener func() (keyring.Keyring, error)) (cleanup func()) {
	original := ringOpener
	ringOpener = opener
	return func() { ringOpener = original }
}

func productionRingOpener() (keyring.Keyring, error) {
	homeDir, _ := os.UserHomeDir()
	return keyring.Open(keyring.Config{
		ServiceName: appName,
		AllowedBackends: []keyring.BackendType{
			keyring.KeychainBackend,
			keyring.SecretServiceBackend,
			keyring.WinCredBackend,
			keyring.FileBackend,
		},
		FileDir:          filepath.Join(homeDir, ".config", "relaxtech", "keyring"),
		FilePasswordFunc: fallbackPassword,
	})
}

func fallbackPassword(_ string) (string, error) {
	return appName + "-fallback", nil
}

// Persist stores a token in the system keyring, associated with the
// given authentication context.
func Persist(ctx AuthContext, tok string, expiry time.Time, uid string) error {
	kr, err := ringOpener()
	if err != nil {
		log.Debug("could not open keyring for saving", "error", err)
		return err
	}
	payload, err := json.Marshal(StoredToken{
		Token:     tok,
		ExpiresAt: expiry,
		UserID:    uid,
	})
	if err != nil {
		return err
	}
	if err := kr.Set(keyring.Item{
		Key:   buildKey(ctx),
		Label: appName + " token",
		Data:  payload,
	}); err != nil {
		log.Debug("failed to write token to keyring", "error", err)
		return err
	}
	log.Debug("token persisted to keyring")
	return nil
}

// Retrieve fetches a cached token for the given context. If the token
// has expired it is automatically removed and ErrKeyNotFound is returned.
// When requiredUID is non-empty, a mismatch with the stored user ID
// causes the lookup to fail.
func Retrieve(ctx AuthContext, requiredUID string) (*StoredToken, error) {
	kr, err := ringOpener()
	if err != nil {
		log.Debug("could not open keyring for retrieval", "error", err)
		return nil, err
	}
	lookupKey := buildKey(ctx)
	entry, err := kr.Get(lookupKey)
	if err == keyring.ErrKeyNotFound && ctx.Email == "" {
		// When no email is provided, try to locate a unique token for
		// this base URL and client ID combination.
		if resolved, resolveErr := locateUniqueEntry(kr, ctx); resolveErr == nil {
			lookupKey = resolved
			entry, err = kr.Get(lookupKey)
		} else {
			log.Debug("broad token search unsuccessful", "error", resolveErr)
		}
	}
	if err != nil {
		log.Debug("token retrieval from keyring failed", "error", err)
		return nil, err
	}
	var stored StoredToken
	if err := json.Unmarshal(entry.Data, &stored); err != nil {
		return nil, err
	}
	// Evict expired tokens automatically.
	if time.Now().After(stored.ExpiresAt) {
		_ = kr.Remove(lookupKey)
		return nil, keyring.ErrKeyNotFound
	}
	// Reject if the stored user does not match the expected one.
	if requiredUID != "" && stored.UserID != "" && stored.UserID != requiredUID {
		return nil, keyring.ErrKeyNotFound
	}
	return &stored, nil
}

// Remove deletes the cached token for the given context. Missing entries
// are silently ignored.
func Remove(ctx AuthContext) error {
	kr, err := ringOpener()
	if err != nil {
		return err
	}
	entryKey := buildKey(ctx)
	if err := kr.Remove(entryKey); err != nil {
		if err == keyring.ErrKeyNotFound || os.IsNotExist(err) {
			return nil
		}
		return err
	}
	return nil
}

// buildKey produces a deterministic, normalized cache key from an AuthContext.
func buildKey(ctx AuthContext) string {
	normalizedURL := strings.TrimSuffix(strings.ToLower(strings.TrimSpace(ctx.BaseURL)), "/")
	normalizedEmail := strings.ToLower(strings.TrimSpace(ctx.Email))
	return storedTokenID + ":" + normalizedURL + "|" + ctx.ClientID + "|" + normalizedEmail
}

// locateUniqueEntry scans all keyring entries for a single match against
// the given base URL and client ID. Returns ErrKeyNotFound when zero or
// more than one candidates exist.
func locateUniqueEntry(kr keyring.Keyring, ctx AuthContext) (string, error) {
	allKeys, err := kr.Keys()
	if err != nil {
		return "", err
	}
	needle := storedTokenID + ":" + strings.TrimSuffix(strings.ToLower(strings.TrimSpace(ctx.BaseURL)), "/") + "|" + ctx.ClientID + "|"
	var candidates []string
	for _, k := range allKeys {
		if strings.HasPrefix(k, needle) {
			candidates = append(candidates, k)
		}
	}
	if len(candidates) == 1 {
		return candidates[0], nil
	}
	return "", keyring.ErrKeyNotFound
}
