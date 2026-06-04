//go:build !integration

package config

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

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
