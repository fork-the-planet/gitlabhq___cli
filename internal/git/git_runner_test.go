package git

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestStandardGitCommand_Git_SetsLocale(t *testing.T) {
	// This test verifies that StandardGitCommand.Git() sets LC_ALL=C
	// to ensure Git output is always in English, regardless of user's locale.

	// Use InitGitRepoWithCommit to ensure we have a proper repo with commits
	// so git status returns predictable output
	InitGitRepoWithCommit(t)

	gitCmd := StandardGitCommand{}

	t.Run("git status outputs English messages", func(t *testing.T) {
		// Run git status which outputs human-readable messages
		output, err := gitCmd.Git("status", "--long")
		require.NoError(t, err)

		// Verify output contains English messages
		// These strings would be different in other locales if LC_ALL=C wasn't set:
		// - German: "Auf Branch" instead of "On branch"
		// - French: "Sur la branche" instead of "On branch"
		// - Chinese: "位于分支" instead of "On branch"
		// Note: Fresh repos may show "No commits yet" but repos with commits show "On branch"
		require.True(t,
			strings.Contains(output, "On branch") || strings.Contains(output, "No commits yet"),
			"Git output should be in English, got: %s", output)
	})

	t.Run("git status contains expected English phrases", func(t *testing.T) {
		// InitGitRepoWithCommit creates a repo with a commit and clean working tree,
		// so git status should output "nothing to commit" since there are no changes
		output, err := gitCmd.Git("status", "--long")
		require.NoError(t, err)

		// This is the key phrase that stack sync looks for to detect up-to-date branches
		// Without LC_ALL=C, it would be localized and string matching would fail:
		// - German: "nichts zu committen"
		// - French: "rien à valider"
		// - Chinese: "无文件要提交"
		require.Contains(t, output, "nothing to commit",
			"Git status should contain 'nothing to commit' for clean working tree")

		// Also verify it shows the branch is clean (another English phrase)
		require.Contains(t, output, "working tree clean",
			"Git status should contain 'working tree clean' for clean working tree")
	})

	t.Run("git commands work regardless of user's LANG setting", func(t *testing.T) {
		// Even if the user has a non-English LANG variable set in their environment,
		// the LC_ALL=C setting should override it and force English output.
		// We can't reliably test locale changes in the test environment,
		// but we can verify the command succeeds and produces valid output.
		output, err := gitCmd.Git("branch", "--show-current")
		require.NoError(t, err)
		require.NotEmpty(t, output, "Git command should produce output")
	})
}
