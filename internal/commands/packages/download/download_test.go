//go:build !integration

package download

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	gitlab "gitlab.com/gitlab-org/api/client-go/v2"
	gitlabtesting "gitlab.com/gitlab-org/api/client-go/v2/testing"

	"gitlab.com/gitlab-org/cli/internal/testing/cmdtest"
)

const (
	fileContents         = "Hello"
	fileContentsChecksum = "185f8db32271fe25f561a6fc938b2e264306ec304eda518007d1764826381969"
	repoName             = "OWNER/REPO"
)

func mockChecksumLookup(tc *gitlabtesting.TestClient, checksum string) {
	tc.MockPackages.EXPECT().
		ListProjectPackages(repoName, gomock.Any(), gomock.Any(), gomock.Any()).
		Return([]*gitlab.Package{{ID: 7, Name: "my-package", Version: "1.0.0"}}, &gitlab.Response{}, nil)
	tc.MockPackages.EXPECT().
		ListPackageFiles(repoName, int64(7), gomock.Any(), gomock.Any(), gomock.Any()).
		Return([]*gitlab.PackageFile{{ID: 1, FileName: "app.zip", FileSHA256: checksum}}, &gitlab.Response{}, nil)
}

func Test_PackagesDownload(t *testing.T) {
	testCases := []struct {
		name             string
		cli              string
		expectedMsg      []string
		expectedFileName string
		wantErr          bool
		wantStderr       string
		setup            func(t *testing.T)
		setupMocks       func(*gitlabtesting.TestClient)
	}{
		{
			name:             "Download with checksum verification",
			cli:              "--name my-package --version 1.0.0 --filename app.zip",
			expectedMsg:      []string{"Downloaded package file to 'app.zip' (Package: my-package, Version: 1.0.0)\n"},
			expectedFileName: "app.zip",
			setupMocks: func(tc *gitlabtesting.TestClient) {
				tc.MockGenericPackages.EXPECT().
					DownloadPackageFile(repoName, "my-package", "1.0.0", "app.zip", gomock.Any()).
					Return([]byte(fileContents), nil, nil)
				mockChecksumLookup(tc, fileContentsChecksum)
			},
		},
		{
			name:             "Download to a custom path",
			cli:              "--name my-package --version 1.0.0 --filename app.zip --path newdir/new.zip",
			expectedMsg:      []string{"Downloaded package file to 'newdir/new.zip' (Package: my-package, Version: 1.0.0)\n"},
			expectedFileName: "newdir/new.zip",
			setupMocks: func(tc *gitlabtesting.TestClient) {
				tc.MockGenericPackages.EXPECT().
					DownloadPackageFile(repoName, "my-package", "1.0.0", "app.zip", gomock.Any()).
					Return([]byte(fileContents), nil, nil)
				mockChecksumLookup(tc, fileContentsChecksum)
			},
		},
		{
			name:             "Download without checksum verification",
			cli:              "--name my-package --version 1.0.0 --filename app.zip --no-verify",
			expectedMsg:      []string{"Downloaded package file to 'app.zip' (Package: my-package, Version: 1.0.0)\n"},
			expectedFileName: "app.zip",
			setupMocks: func(tc *gitlabtesting.TestClient) {
				tc.MockGenericPackages.EXPECT().
					DownloadPackageFile(repoName, "my-package", "1.0.0", "app.zip", gomock.Any()).
					Return([]byte(fileContents), nil, nil)
			},
		},
		{
			name:       "Download with invalid checksum",
			cli:        "--name my-package --version 1.0.0 --filename app.zip",
			wantErr:    true,
			wantStderr: "checksum verification failed for app.zip: expected invalid_checksum, got " + fileContentsChecksum,
			setupMocks: func(tc *gitlabtesting.TestClient) {
				tc.MockGenericPackages.EXPECT().
					DownloadPackageFile(repoName, "my-package", "1.0.0", "app.zip", gomock.Any()).
					Return([]byte(fileContents), nil, nil)
				mockChecksumLookup(tc, "invalid_checksum")
			},
		},
		{
			name:             "Download into a directory keeps the original name",
			cli:              "--name my-package --version 1.0.0 --filename app.zip --path subdir/",
			expectedMsg:      []string{"Downloaded package file to 'subdir/app.zip' (Package: my-package, Version: 1.0.0)\n"},
			expectedFileName: "subdir/app.zip",
			setupMocks: func(tc *gitlabtesting.TestClient) {
				tc.MockGenericPackages.EXPECT().
					DownloadPackageFile(repoName, "my-package", "1.0.0", "app.zip", gomock.Any()).
					Return([]byte(fileContents), nil, nil)
				mockChecksumLookup(tc, fileContentsChecksum)
			},
		},
		{
			name:             "Download into an existing directory keeps the original name",
			cli:              "--name my-package --version 1.0.0 --filename app.zip --path existing",
			expectedMsg:      []string{"Downloaded package file to 'existing/app.zip' (Package: my-package, Version: 1.0.0)\n"},
			expectedFileName: "existing/app.zip",
			setup: func(t *testing.T) {
				t.Helper()
				require.NoError(t, os.Mkdir("existing", 0o755))
			},
			setupMocks: func(tc *gitlabtesting.TestClient) {
				tc.MockGenericPackages.EXPECT().
					DownloadPackageFile(repoName, "my-package", "1.0.0", "app.zip", gomock.Any()).
					Return([]byte(fileContents), nil, nil)
				mockChecksumLookup(tc, fileContentsChecksum)
			},
		},
		{
			name:             "Force overwrites an existing file",
			cli:              "--name my-package --version 1.0.0 --filename app.zip --force",
			expectedMsg:      []string{"Downloaded package file to 'app.zip' (Package: my-package, Version: 1.0.0)\n"},
			expectedFileName: "app.zip",
			setup: func(t *testing.T) {
				t.Helper()
				require.NoError(t, os.WriteFile("app.zip", []byte("stale"), 0o644))
			},
			setupMocks: func(tc *gitlabtesting.TestClient) {
				tc.MockGenericPackages.EXPECT().
					DownloadPackageFile(repoName, "my-package", "1.0.0", "app.zip", gomock.Any()).
					Return([]byte(fileContents), nil, nil)
				mockChecksumLookup(tc, fileContentsChecksum)
			},
		},
		{
			name:       "Existing file without force errors",
			cli:        "--name my-package --version 1.0.0 --filename app.zip",
			wantErr:    true,
			wantStderr: "file app.zip already exists; use --force to overwrite current file",
			setup: func(t *testing.T) {
				t.Helper()
				require.NoError(t, os.WriteFile("app.zip", []byte("stale"), 0o644))
			},
			setupMocks: func(tc *gitlabtesting.TestClient) {},
		},
		{
			name:       "Download but API errors",
			cli:        "--name my-package --version 1.0.0 --filename app.zip --no-verify",
			wantErr:    true,
			wantStderr: "failed to download package file: GET https://gitlab.com/api/v4/projects/OWNER%2FREPO/packages/generic: 404",
			setupMocks: func(tc *gitlabtesting.TestClient) {
				tc.MockGenericPackages.EXPECT().
					DownloadPackageFile(repoName, "my-package", "1.0.0", "app.zip", gomock.Any()).
					Return(nil, nil, fmt.Errorf("GET https://gitlab.com/api/v4/projects/OWNER%%2FREPO/packages/generic: 404"))
			},
		},
		{
			name:       "Download missing file in package",
			cli:        "--name my-package --version 1.0.0 --filename missing.zip",
			wantErr:    true,
			wantStderr: "couldn't locate file missing.zip in package my-package version 1.0.0",
			setupMocks: func(tc *gitlabtesting.TestClient) {
				tc.MockGenericPackages.EXPECT().
					DownloadPackageFile(repoName, "my-package", "1.0.0", "missing.zip", gomock.Any()).
					Return([]byte(fileContents), nil, nil)
				mockChecksumLookup(tc, fileContentsChecksum)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tempDir := t.TempDir()
			t.Chdir(tempDir)

			if tc.setup != nil {
				tc.setup(t)
			}

			testClient := gitlabtesting.NewTestClient(t)
			tc.setupMocks(testClient)

			exec := cmdtest.SetupCmdForTest(t, NewCmd, false,
				cmdtest.WithGitLabClient(testClient.Client),
			)
			out, err := exec(tc.cli)

			if tc.wantErr {
				require.Error(t, err)
				assert.Equal(t, tc.wantStderr, err.Error())
				return
			}
			require.NoError(t, err)

			for _, msg := range tc.expectedMsg {
				require.Contains(t, out.String(), msg)
			}

			actualContent, err := os.ReadFile(tc.expectedFileName)
			require.NoError(t, err, "Failed to read downloaded test file")
			assert.Equal(t, fileContents, string(actualContent), "File content should match")
		})
	}
}
