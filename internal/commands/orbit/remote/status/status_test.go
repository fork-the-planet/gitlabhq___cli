//go:build !integration

package status

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	gitlab "gitlab.com/gitlab-org/api/client-go/v2"
	gitlabtesting "gitlab.com/gitlab-org/api/client-go/v2/testing"

	"gitlab.com/gitlab-org/cli/internal/api"
	"gitlab.com/gitlab-org/cli/internal/cmdutils"
	"gitlab.com/gitlab-org/cli/internal/commands/orbit/internal/orbiterr"
	"gitlab.com/gitlab-org/cli/internal/testing/cmdtest"
)

func TestStatus_NewShape_HappyPath(t *testing.T) {
	t.Parallel()
	// GIVEN the new-shape API returns user.available=true with system health
	testClient := gitlabtesting.NewTestClient(t)
	testClient.MockOrbit.EXPECT().
		GetStatus(gomock.Any(), gomock.Any()).
		Return(&gitlab.OrbitStatus{
			User: &gitlab.OrbitStatusUser{Available: true},
			System: &gitlab.OrbitStatusSystem{
				Status:    "healthy",
				Timestamp: "2026-06-20T12:00:00Z",
				Version:   "0.6.0",
				Components: []*gitlab.OrbitStatusComponent{
					{
						Name:     "clickhouse",
						Status:   "healthy",
						Replicas: &gitlab.OrbitStatusReplicas{Ready: 3, Desired: 3},
					},
				},
			},
			// Flat fields omitted: run() reads only status.System
			// for the new shape. UnmarshalJSON promotes them in
			// real usage, but they are not consulted here.
		}, &gitlab.Response{Response: &http.Response{StatusCode: http.StatusOK}}, nil)

	exec := cmdtest.SetupCmdForTest(
		t,
		NewCmd,
		false,
		cmdtest.WithApiClient(cmdtest.NewTestApiClient(t, nil, "", "", api.WithGitLabClient(testClient.Client))),
	)

	// WHEN `glab orbit remote status` runs
	out, err := exec("")

	// THEN the system health is printed as JSON (not the full wrapper)
	require.NoError(t, err)
	assert.Empty(t, out.ErrBuf.String())

	var result gitlab.OrbitStatusSystem
	require.NoError(t, json.Unmarshal(out.OutBuf.Bytes(), &result))
	assert.Equal(t, "healthy", result.Status)
	assert.Equal(t, "0.6.0", result.Version)
	require.Len(t, result.Components, 1)
	assert.Equal(t, "clickhouse", result.Components[0].Name)
}

func TestStatus_NewShape_UserUnavailable(t *testing.T) {
	t.Parallel()
	// GIVEN the new-shape API returns user.available=false (no access)
	testClient := gitlabtesting.NewTestClient(t)
	testClient.MockOrbit.EXPECT().
		GetStatus(gomock.Any(), gomock.Any()).
		Return(&gitlab.OrbitStatus{
			User:   &gitlab.OrbitStatusUser{Available: false},
			System: nil,
		}, &gitlab.Response{Response: &http.Response{StatusCode: http.StatusOK}}, nil)

	exec := cmdtest.SetupCmdForTest(
		t,
		NewCmd,
		false,
		cmdtest.WithApiClient(cmdtest.NewTestApiClient(t, nil, "", "", api.WithGitLabClient(testClient.Client))),
	)

	// WHEN `glab orbit remote status` runs
	_, err := exec("")

	// THEN the error maps to ExitOrbitUnavailable (exit code 2)
	require.Error(t, err)
	var exitErr *cmdutils.ExitError
	require.ErrorAs(t, err, &exitErr)
	assert.Equal(t, orbiterr.ExitOrbitUnavailable, exitErr.Code)
	assert.Contains(t, err.Error(), "not available for your user")
}

func TestStatus_NewShape_UserAvailable_SystemNil(t *testing.T) {
	t.Parallel()
	// GIVEN the new-shape API returns user.available=true but no system object
	testClient := gitlabtesting.NewTestClient(t)
	testClient.MockOrbit.EXPECT().
		GetStatus(gomock.Any(), gomock.Any()).
		Return(&gitlab.OrbitStatus{
			User:   &gitlab.OrbitStatusUser{Available: true},
			System: nil,
		}, &gitlab.Response{Response: &http.Response{StatusCode: http.StatusOK}}, nil)

	exec := cmdtest.SetupCmdForTest(
		t,
		NewCmd,
		false,
		cmdtest.WithApiClient(cmdtest.NewTestApiClient(t, nil, "", "", api.WithGitLabClient(testClient.Client))),
	)

	// WHEN `glab orbit remote status` runs
	_, err := exec("")

	// THEN an explicit error is returned (not the user/system wrapper)
	require.Error(t, err)
	var exitErr *cmdutils.ExitError
	require.ErrorAs(t, err, &exitErr)
	assert.Equal(t, orbiterr.ExitOrbitUnavailable, exitErr.Code)
	assert.Contains(t, err.Error(), "system health is absent")
}

func TestStatus_OldFlatShape(t *testing.T) {
	t.Parallel()
	// GIVEN an older instance returns the flat health shape (no User/System)
	testClient := gitlabtesting.NewTestClient(t)
	testClient.MockOrbit.EXPECT().
		GetStatus(gomock.Any(), gomock.Any()).
		Return(&gitlab.OrbitStatus{
			Status:    "healthy",
			Timestamp: "2026-04-28T12:00:00Z",
			Version:   "0.5.0",
			Components: []*gitlab.OrbitStatusComponent{
				{
					Name:     "clickhouse",
					Status:   "healthy",
					Replicas: &gitlab.OrbitStatusReplicas{Ready: 3, Desired: 3},
				},
			},
		}, &gitlab.Response{Response: &http.Response{StatusCode: http.StatusOK}}, nil)

	exec := cmdtest.SetupCmdForTest(
		t,
		NewCmd,
		false,
		cmdtest.WithApiClient(cmdtest.NewTestApiClient(t, nil, "", "", api.WithGitLabClient(testClient.Client))),
	)

	// WHEN `glab orbit remote status` runs
	out, err := exec("")

	// THEN the flat status response is printed as JSON (back-compat)
	require.NoError(t, err)
	assert.Empty(t, out.ErrBuf.String())

	var result gitlab.OrbitStatus
	require.NoError(t, json.Unmarshal(out.OutBuf.Bytes(), &result))
	assert.Equal(t, "healthy", result.Status)
	assert.Equal(t, "0.5.0", result.Version)
	require.Len(t, result.Components, 1)
	assert.Equal(t, "clickhouse", result.Components[0].Name)
}

func TestStatus_FeatureFlagOff(t *testing.T) {
	t.Parallel()
	// GIVEN the API returns 404 because the knowledge_graph FF is off
	testClient := gitlabtesting.NewTestClient(t)
	testClient.MockOrbit.EXPECT().
		GetStatus(gomock.Any(), gomock.Any()).
		Return(nil,
			&gitlab.Response{Response: &http.Response{StatusCode: http.StatusNotFound}},
			&gitlab.ErrorResponse{
				Response: &http.Response{StatusCode: http.StatusNotFound},
				Message:  "404 Not Found",
			})

	exec := cmdtest.SetupCmdForTest(
		t,
		NewCmd,
		false,
		cmdtest.WithApiClient(cmdtest.NewTestApiClient(t, nil, "", "", api.WithGitLabClient(testClient.Client))),
	)

	// WHEN `glab orbit remote status` runs
	_, err := exec("")

	// THEN the error is mapped to ExitOrbitUnavailable (exit code 2)
	require.Error(t, err)
	var exitErr *cmdutils.ExitError
	require.ErrorAs(t, err, &exitErr)
	assert.Equal(t, orbiterr.ExitOrbitUnavailable, exitErr.Code)
}
