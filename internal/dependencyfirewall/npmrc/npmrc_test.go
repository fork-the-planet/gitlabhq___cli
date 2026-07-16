//go:build !integration

package npmrc

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"gitlab.com/gitlab-org/cli/internal/dependencyfirewall/config"
)

func lines(body string) []string {
	return strings.Split(strings.TrimRight(body, "\n"), "\n")
}

func settings() config.Settings {
	return config.Settings{
		RegistryURL: "https://gitlab.example.com/api/v4/projects/g%2Fp/packages/npm/",
		AuthHost:    "gitlab.example.com/api/v4/projects/g%2Fp/packages/npm/",
		AuthToken:   "tok123",
	}
}

func TestApplyWritesFirewallNpmrcWhenNoneExists(t *testing.T) {
	dir := t.TempDir()
	h, err := Apply(dir, settings())
	require.NoError(t, err)
	defer func() { _ = h.Restore() }()

	raw, err := os.ReadFile(filepath.Join(dir, ".npmrc"))
	require.NoError(t, err)
	body := string(raw)
	assert.Contains(t, body, "registry=https://gitlab.example.com/api/v4/projects/g%2Fp/packages/npm/")
	assert.Contains(t, body, "//gitlab.example.com/api/v4/projects/g%2Fp/packages/npm/:_authToken=tok123")
	assert.NotContains(t, body, "strict-ssl=")
}

func TestApplyPreservesUserStrictSsl(t *testing.T) {
	dir := t.TempDir()
	original := "strict-ssl=true\nfoo=bar\n"
	require.NoError(t, os.WriteFile(filepath.Join(dir, ".npmrc"), []byte(original), 0o644))

	h, err := Apply(dir, config.Settings{
		RegistryURL: "https://gitlab.example.com/api/v4/projects/42/packages/npm/",
		AuthHost:    "gitlab.example.com/api/v4/projects/42/packages/npm/",
		AuthToken:   "tok123",
	})
	require.NoError(t, err)
	defer func() { _ = h.Restore() }()

	raw, err := os.ReadFile(filepath.Join(dir, ".npmrc"))
	require.NoError(t, err)
	body := string(raw)
	assert.Contains(t, body, "strict-ssl=true")
	assert.NotContains(t, body, "strict-ssl=false")
	assert.Contains(t, body, "registry=https://gitlab.example.com/api/v4/projects/42/packages/npm/")
}

func TestApplyPreservesScopedRegistries(t *testing.T) {
	dir := t.TempDir()
	original := "registry=https://registry.npmjs.org/\n@acme:registry=https://npm.acme.test/\n"
	require.NoError(t, os.WriteFile(filepath.Join(dir, ".npmrc"), []byte(original), 0o644))

	h, err := Apply(dir, settings())
	require.NoError(t, err)
	defer func() { _ = h.Restore() }()

	raw, err := os.ReadFile(filepath.Join(dir, ".npmrc"))
	require.NoError(t, err)
	body := string(raw)
	assert.NotContains(t, lines(body), "registry=https://registry.npmjs.org/")
	assert.Contains(t, body, "@acme:registry=https://npm.acme.test/")
	assert.Contains(t, body, settings().RegistryURL)
}

func TestRestoreDeletesWrittenNpmrcWhenNoOriginal(t *testing.T) {
	dir := t.TempDir()
	h, err := Apply(dir, settings())
	require.NoError(t, err)
	require.NoError(t, h.Restore())

	_, err = os.Stat(filepath.Join(dir, ".npmrc"))
	assert.True(t, os.IsNotExist(err), ".npmrc should be removed")
	_, err = os.Stat(filepath.Join(dir, ".gitlab.npmrc.backup"))
	assert.True(t, os.IsNotExist(err), "backup should be removed")
}

func TestApplyAndRestoreLeavesOriginalByteIdentical(t *testing.T) {
	dir := t.TempDir()
	original := "registry=https://registry.npmjs.org/\n@acme:registry=https://registry.npmjs.org/\nfoo=bar\n"
	require.NoError(t, os.WriteFile(filepath.Join(dir, ".npmrc"), []byte(original), 0o644))

	h, err := Apply(dir, settings())
	require.NoError(t, err)

	rewritten, err := os.ReadFile(filepath.Join(dir, ".npmrc"))
	require.NoError(t, err)
	assert.NotContains(t, lines(string(rewritten)), "registry=https://registry.npmjs.org/")
	assert.Contains(t, string(rewritten), "@acme:registry=https://registry.npmjs.org/")
	assert.Contains(t, string(rewritten), "foo=bar")

	require.NoError(t, h.Restore())

	after, err := os.ReadFile(filepath.Join(dir, ".npmrc"))
	require.NoError(t, err)
	assert.Equal(t, original, string(after), ".npmrc must be byte-identical after restore")

	_, err = os.Stat(filepath.Join(dir, ".gitlab.npmrc.backup"))
	assert.True(t, os.IsNotExist(err), "backup must be deleted")
}

func assertOwnerOnly(t *testing.T, path string) {
	t.Helper()
	info, err := os.Stat(path)
	require.NoError(t, err)
	assert.Equal(t, os.FileMode(0o600), info.Mode().Perm(),
		"%s may contain an auth token and must be owner read/write only", path)
}

func TestApplyWritesTokenFilesOwnerOnly(t *testing.T) {
	dir := t.TempDir()
	require.NoError(t, os.WriteFile(filepath.Join(dir, ".npmrc"), []byte("foo=bar\n"), 0o644))

	h, err := Apply(dir, settings())
	require.NoError(t, err)
	defer func() { _ = h.Restore() }()

	assertOwnerOnly(t, filepath.Join(dir, ".npmrc"))
	assertOwnerOnly(t, filepath.Join(dir, ".gitlab.npmrc.backup"))
}

func TestRestoreWritesNpmrcOwnerOnly(t *testing.T) {
	dir := t.TempDir()
	require.NoError(t, os.WriteFile(filepath.Join(dir, ".npmrc"), []byte("foo=bar\n"), 0o644))

	h, err := Apply(dir, settings())
	require.NoError(t, err)
	require.NoError(t, h.Restore())

	assertOwnerOnly(t, filepath.Join(dir, ".npmrc"))
}

func TestApplyRefusesWhenBackupAlreadyExists(t *testing.T) {
	dir := t.TempDir()
	original := "foo=bar\n"
	require.NoError(t, os.WriteFile(filepath.Join(dir, ".npmrc"), []byte(original), 0o644))

	h, err := Apply(dir, settings())
	require.NoError(t, err)
	_ = h

	_, err2 := Apply(dir, settings())
	require.Error(t, err2, "second Apply must refuse while a backup exists")

	backup, err := os.ReadFile(filepath.Join(dir, ".gitlab.npmrc.backup"))
	require.NoError(t, err)
	assert.Equal(t, original, string(backup),
		"backup must still hold the pristine original, not the modified .npmrc containing the auth token")
	assert.NotContains(t, string(backup), "_authToken=",
		"backup must never contain an auth token")

	require.NoError(t, h.Restore())
	after, err := os.ReadFile(filepath.Join(dir, ".npmrc"))
	require.NoError(t, err)
	assert.Equal(t, original, string(after), "original must survive a refused second Apply")
}

func TestApplyOmitsAuthLineWhenNoToken(t *testing.T) {
	dir := t.TempDir()
	s := config.Settings{
		RegistryURL: "https://registry.npmjs.org/",
		AuthHost:    "registry.npmjs.org/",
		AuthToken:   "",
	}
	h, err := Apply(dir, s)
	require.NoError(t, err)
	defer func() { _ = h.Restore() }()

	raw, err := os.ReadFile(filepath.Join(dir, ".npmrc"))
	require.NoError(t, err)
	body := string(raw)
	assert.Contains(t, body, "registry=https://registry.npmjs.org/")
	assert.NotContains(t, body, "_authToken=")
}
