package helpers

import (
	"fmt"

	gitlab "gitlab.com/gitlab-org/api/client-go/v2"

	"gitlab.com/gitlab-org/cli/internal/api"
)

func GetSecureFileIDByName(client *gitlab.Client, fileName, repoName string) (int64, error) {
	options := &gitlab.ListProjectSecureFilesOptions{
		ListOptions: gitlab.ListOptions{
			Page:    1,
			PerPage: api.MaxPerPage,
		},
	}

	for secureFile, err := range gitlab.Scan2(func(p gitlab.PaginationOptionFunc) ([]*gitlab.SecureFile, *gitlab.Response, error) {
		return client.SecureFiles.ListProjectSecureFiles(repoName, options, p)
	}) {
		if err != nil {
			return 0, fmt.Errorf("error fetching secure files: %w", err)
		}

		if secureFile.Name == fileName {
			return secureFile.ID, nil
		}
	}

	return 0, fmt.Errorf("couldn't locate secure file with name %s", fileName)
}
