//go:build !integration

package incident

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"gitlab.com/gitlab-org/cli/internal/testing/cmdtest"
)

func TestIncidentCmd(t *testing.T) {
	var buf bytes.Buffer
	cmd := NewCmdIncident(cmdtest.NewTestFactory(nil))
	cmd.SetOut(&buf)

	require.NoError(t, cmd.Execute())

	assert.Contains(t, buf.String(), "Incidents track service disruptions")
}
