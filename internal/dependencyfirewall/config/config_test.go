//go:build !integration

package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoadMissingReturnsEmpty(t *testing.T) {
	c, err := Load(t.TempDir())
	require.NoError(t, err)
	assert.Empty(t, c.NPM.RepoResolve)
	assert.Empty(t, c.NPM.RepoDeploy)
}

func TestMergeThenLoadRoundTrip(t *testing.T) {
	dir := t.TempDir()
	resolve, deploy := "virt", "local"
	require.NoError(t, Merge(dir, "npm", &resolve, &deploy))

	got, err := Load(dir)
	require.NoError(t, err)
	assert.Equal(t, "virt", got.NPM.RepoResolve)
	assert.Equal(t, "local", got.NPM.RepoDeploy)
}

func TestMergePreservesUnknownKeysAndUpdatesOnlySetFields(t *testing.T) {
	dir := t.TempDir()
	require.NoError(t, os.MkdirAll(filepath.Join(dir, ".gitlab", "df"), 0o755))
	seed := `{"npm":{"repoResolve":"oldR","repoDeploy":"oldD"},"future":{"x":1}}`
	require.NoError(t, os.WriteFile(filepath.Join(dir, ".gitlab", "df", "config.json"), []byte(seed), 0o644))

	resolve := "newR"
	require.NoError(t, Merge(dir, "npm", &resolve, nil))

	raw, err := os.ReadFile(filepath.Join(dir, ".gitlab", "df", "config.json"))
	require.NoError(t, err)
	assert.Contains(t, string(raw), `"repoResolve": "newR"`)
	assert.Contains(t, string(raw), `"repoDeploy": "oldD"`)
	assert.Contains(t, string(raw), `"future"`)
}

func TestMergeRejectsNonObjectNpmKey(t *testing.T) {
	dir := t.TempDir()
	require.NoError(t, os.MkdirAll(filepath.Join(dir, ".gitlab", "df"), 0o755))
	require.NoError(t, os.WriteFile(filepath.Join(dir, ".gitlab", "df", "config.json"),
		[]byte(`{"npm":"not-an-object"}`), 0o644))
	resolve := "x"
	err := Merge(dir, "npm", &resolve, nil)
	require.Error(t, err)
}

func TestMergeWithoutDeployWritesNoDeployKey(t *testing.T) {
	dir := t.TempDir()
	require.NoError(t, os.MkdirAll(filepath.Join(dir, ".gitlab", "df"), 0o755))
	resolve := "https://npm.example/"
	require.NoError(t, Merge(dir, "npm", &resolve, nil))

	raw, err := os.ReadFile(filepath.Join(dir, ".gitlab", "df", "config.json"))
	require.NoError(t, err)
	assert.Contains(t, string(raw), `"npm"`)
	assert.Contains(t, string(raw), `"repoResolve": "https://npm.example/"`)
	assert.NotContains(t, string(raw), `"repoDeploy"`)
}

func TestMergePreservesOtherManagerBlocks(t *testing.T) {
	dir := t.TempDir()
	require.NoError(t, os.MkdirAll(filepath.Join(dir, ".gitlab", "df"), 0o755))
	seed := `{"other":{"repoResolve":"otherR"}}`
	require.NoError(t, os.WriteFile(filepath.Join(dir, ".gitlab", "df", "config.json"), []byte(seed), 0o644))

	resolve := "npmR"
	require.NoError(t, Merge(dir, "npm", &resolve, nil))

	raw, err := os.ReadFile(filepath.Join(dir, ".gitlab", "df", "config.json"))
	require.NoError(t, err)
	assert.Contains(t, string(raw), `"otherR"`)
	assert.Contains(t, string(raw), `"npmR"`)
}
