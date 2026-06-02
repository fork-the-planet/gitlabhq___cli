//go:build !integration

package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestKnownKeys_IncludesUserSettableKeysFromSchema(t *testing.T) {
	t.Parallel()

	keys := KnownKeys()

	for _, k := range []string{
		// Global scalars.
		"git_protocol",
		"editor",
		"browser",
		"glamour_style",
		"host",
		"no_prompt",
		"telemetry",
		"branch_prefix",
		"remote_alias",
		// Per-host scalars.
		"token",
		"job_token",
		"user",
		"client_id",
		"ca_cert",
		"client_cert",
		"client_key",
		"skip_tls_verify",
		"use_keyring",
		"api_protocol",
		"api_host",
	} {
		_, ok := keys[k]
		assert.True(t, ok, "expected %q to be a known key", k)
	}
}

func TestKnownKeys_ExcludesInternalStateKeys(t *testing.T) {
	t.Parallel()

	keys := KnownKeys()

	// Keys flagged UserSettable: false in the schema must not appear.
	for _, k := range []string{
		"is_oauth2",
		"oauth2_refresh_token",
		"oauth2_expiry_date",
		"refresh_token",
	} {
		_, ok := keys[k]
		assert.False(t, ok, "expected internal-state key %q to be excluded", k)
	}

	// Structural containers must not appear.
	for _, k := range []string{"hosts", "aliases"} {
		_, ok := keys[k]
		assert.False(t, ok, "expected structural key %q to be excluded", k)
	}
}

func TestIsKnownKey(t *testing.T) {
	t.Parallel()

	assert.True(t, IsKnownKey("editor"))
	assert.True(t, IsKnownKey("token"))
	assert.False(t, IsKnownKey("oauth_scopes"))
	assert.False(t, IsKnownKey("oauth2_refresh_token"), "internal-state keys are not user-settable")
	assert.False(t, IsKnownKey(""))
	assert.False(t, IsKnownKey("hosts"))
}
