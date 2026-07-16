//go:build !integration

package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRegistrySettingsProjectURLAttachesTokenForGitLabHost(t *testing.T) {
	s, err := RegistrySettings(
		"gitlab.example.com",
		"https://gitlab.example.com/api/v4/projects/my-group%2Fmy-virt/packages/npm/",
		"tok123",
	)
	require.NoError(t, err)
	assert.Equal(t,
		"https://gitlab.example.com/api/v4/projects/my-group%2Fmy-virt/packages/npm/",
		s.RegistryURL)
	assert.Equal(t,
		"gitlab.example.com/api/v4/projects/my-group%2Fmy-virt/packages/npm/",
		s.AuthHost)
	assert.Equal(t, "tok123", s.AuthToken)
}

func TestRegistrySettingsGroupURL(t *testing.T) {
	s, err := RegistrySettings(
		"gitlab.example.com",
		"https://gitlab.example.com/api/v4/groups/42/-/packages/npm/",
		"tok123",
	)
	require.NoError(t, err)
	assert.Equal(t,
		"https://gitlab.example.com/api/v4/groups/42/-/packages/npm/",
		s.RegistryURL)
	assert.Equal(t,
		"gitlab.example.com/api/v4/groups/42/-/packages/npm/",
		s.AuthHost)
	assert.Equal(t, "tok123", s.AuthToken)
}

func TestRegistrySettingsNormalizesTrailingSlash(t *testing.T) {
	s, err := RegistrySettings(
		"gitlab.example.com",
		"https://gitlab.example.com/api/v4/groups/42/-/packages/npm",
		"tok123",
	)
	require.NoError(t, err)
	assert.Equal(t,
		"https://gitlab.example.com/api/v4/groups/42/-/packages/npm/",
		s.RegistryURL)
	assert.Equal(t,
		"gitlab.example.com/api/v4/groups/42/-/packages/npm/",
		s.AuthHost)
}

func TestRegistrySettingsPublicRegistryOmitsToken(t *testing.T) {
	s, err := RegistrySettings(
		"gitlab.example.com",
		"https://registry.npmjs.org/",
		"tok123",
	)
	require.NoError(t, err)
	assert.Equal(t, "https://registry.npmjs.org/", s.RegistryURL)
	assert.Equal(t, "registry.npmjs.org/", s.AuthHost)
	assert.Empty(t, s.AuthToken, "token must not leak to a non-GitLab host")
}

func TestRegistrySettingsEmptyTokenOmitsToken(t *testing.T) {
	s, err := RegistrySettings(
		"gitlab.example.com",
		"https://gitlab.example.com/api/v4/groups/42/-/packages/npm/",
		"",
	)
	require.NoError(t, err)
	assert.Empty(t, s.AuthToken)
}

func TestRegistrySettingsRequiresResolveURL(t *testing.T) {
	_, err := RegistrySettings("gitlab.example.com", "", "tok")
	require.Error(t, err)
}

func TestRegistrySettingsRejectsNonAbsoluteURL(t *testing.T) {
	_, err := RegistrySettings("gitlab.example.com", "my-group/my-virt", "tok")
	require.Error(t, err)
}

func TestRegistrySettingsRejectsNonHTTPURL(t *testing.T) {
	_, err := RegistrySettings("gitlab.example.com", "ftp://example.com/npm/", "tok")
	require.Error(t, err)
}

func TestRegistrySettingsHostMatchIsCaseInsensitive(t *testing.T) {
	s, err := RegistrySettings(
		"gitlab.example.com",
		"https://GitLab.Example.com/api/v4/groups/42/-/packages/npm/",
		"tok123",
	)
	require.NoError(t, err)
	assert.Equal(t, "tok123", s.AuthToken)
}

func TestRegistrySettingsHostMatchIgnoresDefaultPort(t *testing.T) {
	s, err := RegistrySettings(
		"gitlab.example.com",
		"https://gitlab.example.com:443/api/v4/groups/42/-/packages/npm/",
		"tok123",
	)
	require.NoError(t, err)
	assert.Equal(t, "tok123", s.AuthToken)
}

func TestRegistrySettingsHostMatchRejectsNonDefaultPort(t *testing.T) {
	s, err := RegistrySettings(
		"gitlab.example.com",
		"https://gitlab.example.com:8443/api/v4/groups/42/-/packages/npm/",
		"tok123",
	)
	require.NoError(t, err)
	assert.Empty(t, s.AuthToken, "non-default port is a different host; token must not be attached")
}
