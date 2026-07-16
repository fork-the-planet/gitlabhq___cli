//go:build !windows

package fsx

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/google/renameio/v2"
)

// WriteOwnerOnly writes data to path atomically and forces mode 0o600.
//
// On POSIX the write is a temp-file-plus-rename via renameio, so a reader
// concurrent with the write either sees the full old file or the full new
// file — never a half-written mix. If the process dies mid-write, the target
// is untouched and only an orphaned temp file may be left behind.
//
// renameio.WriteFile defaults to WithExistingPermissions(), which preserves
// the existing mode when the target already exists. That is the wrong policy
// for token-bearing files (see the package-level doc comment for the full
// rationale), so we follow the write with an explicit Chmod to guarantee
// 0o600 in both the create and the overwrite paths.
//
// Callers use this for any file that may embed an inline authentication
// token — for example, .npmrc, .yarnrc.yml, or a patched Pipfile — and for
// the token-bearing backups those files produce.
func WriteOwnerOnly(path string, data []byte) error {
	if err := renameio.WriteFile(path, data, 0o600); err != nil {
		return err
	}
	return os.Chmod(path, 0o600)
}

// WriteJSONFile creates path's parent directory (0o755), marshals v with
// two-space indent for user-editability, appends a trailing newline
// (POSIX text-file convention; keeps `git diff` clean when the file is
// hand-edited later), and writes the result to path with mode 0o644.
// The write is atomic via renameio. Use it for non-secret configuration
// and log files; token-bearing files must use WriteOwnerOnly instead.
func WriteJSONFile(path string, v any) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	raw, err := json.MarshalIndent(v, "", "  ") //nolint:forbidigo // writing config/log to disk, not stdout
	if err != nil {
		return err
	}
	return renameio.WriteFile(path, append(raw, '\n'), 0o644)
}
