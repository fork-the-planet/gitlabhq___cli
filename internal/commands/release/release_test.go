//go:build !integration

package release

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"gitlab.com/gitlab-org/cli/internal/testing/cmdtest"
)

func Test_Release(t *testing.T) {
	var buf bytes.Buffer
	cmd := NewCmdRelease(cmdtest.NewTestFactory(nil))
	cmd.SetOut(&buf)

	assert.NotNil(t, cmd.Root())
	require.NoError(t, cmd.Execute())

	assert.Contains(t, buf.String(), "A release bundles a Git tag")
}
