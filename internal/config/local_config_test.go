//go:build !integration

package config

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// An in-memory config (no directory behind it) must not persist local config to
// the surrounding git repository's .git/glab-cli/config.yml, even when Set()
// (which writes) is called from inside a git checkout.
func Test_InMemoryConfig_LocalSetDoesNotPersist(t *testing.T) {
	dir := t.TempDir()
	require.NoError(t, exec.Command("git", "-C", dir, "init").Run())
	t.Chdir(dir)

	cfg := NewBlankConfig()
	local, err := cfg.Local()
	require.NoError(t, err)
	require.NoError(t, local.Set("git_protocol", "ssh"))

	_, statErr := os.Stat(filepath.Join(dir, ".git", "glab-cli", "config.yml"))
	assert.True(t, os.IsNotExist(statErr), "in-memory config must not persist local config to .git")
}

func Test_GitDir(t *testing.T) {
	gotRelative := filepath.Join(GitDir(true)...)
	gotAbsolute := filepath.Join(GitDir(false)...)
	absRelative, err := filepath.Abs(gotRelative)
	require.NoError(t, err)
	assert.Equal(t, gotAbsolute, absRelative)
}

func Test_LocalConfigDir(t *testing.T) {
	got := LocalConfigDir()
	assert.ElementsMatch(t, []string{filepath.Join("..", "..", ".git"), "glab-cli"}, got)
}

func Test_LocalConfigFile(t *testing.T) {
	expectedPath := filepath.Join("..", "..", ".git", "glab-cli", "config.yml")
	got := LocalConfigFile()
	assert.Equal(t, expectedPath, got)
}
