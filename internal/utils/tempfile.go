package utils

import (
	"fmt"
	"math/rand"
	"os"
	"strconv"
)

// CreateTemp is a modified implementation of os.CreateTemp() using root.OpenFile.
func CreateTemp(root *os.Root, path string) (*os.File, error) {
	const retryLimit = 10000

	var (
		f   *os.File
		err error
	)

	// This retry logic is to handle tempfile name collisions with an existing tempfile.
	// This is probably overkill since the chances of a collision are already extremely unlikely.
	// But it is taken from the os.CreateTemp implementation, and makes a collision effectively impossible.
	for range retryLimit {
		path += strconv.Itoa(rand.Intn(10))
		f, err = root.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_EXCL, 0o600)
		if os.IsExist(err) {
			continue
		}

		return f, err
	}

	return nil, fmt.Errorf("failed to create tempfile after %d tries: %w", retryLimit, err)
}
