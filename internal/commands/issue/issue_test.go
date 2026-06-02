//go:build !integration

package issue

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"gitlab.com/gitlab-org/cli/internal/testing/cmdtest"
)

func TestIssueCmd(t *testing.T) {
	var buf bytes.Buffer
	cmd := NewCmdIssue(cmdtest.NewTestFactory(nil))
	cmd.SetOut(&buf)

	require.NoError(t, cmd.Execute())

	assert.Contains(t, buf.String(), "Open issues, list and view them")
}
