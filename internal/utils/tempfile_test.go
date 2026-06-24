package utils

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateTemp(t *testing.T) {
	t.Parallel()

	tmpDir := t.TempDir()
	root, err := os.OpenRoot(tmpDir)
	require.NoError(t, err)

	path := "file.txt"
	f, err := CreateTemp(root, path)
	require.NoError(t, err)
	require.NotNil(t, f)
	got := filepath.Base(f.Name())
	assert.True(t, strings.HasPrefix(got, path))
}
