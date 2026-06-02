//go:build !integration

package board

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"gitlab.com/gitlab-org/cli/internal/testing/cmdtest"
)

func TestNewCmdBoard(t *testing.T) {
	var buf bytes.Buffer
	cmd := NewCmdBoard(cmdtest.NewTestFactory(nil))
	cmd.SetOut(&buf)

	require.NoError(t, cmd.Execute())

	assert.Contains(t, buf.String(), "Issue boards organize issues into lists")
}
