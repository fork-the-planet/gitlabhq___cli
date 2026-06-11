//go:build !integration

package update

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"gitlab.com/gitlab-org/cli/internal/commands/skills/bundled"
	"gitlab.com/gitlab-org/cli/internal/commands/skills/installed"
	"gitlab.com/gitlab-org/cli/internal/commands/skills/skill"
	"gitlab.com/gitlab-org/cli/internal/config"
	"gitlab.com/gitlab-org/cli/internal/testing/cmdtest"
)

func newSkillCheckConfig(t *testing.T, notifySetting string) config.Config {
	t.Helper()
	cfg := config.NewBlankConfig()
	if notifySetting != "" {
		require.NoError(t, cfg.Set("", "notify_skill_updates", notifySetting))
	}
	return cfg
}

// stubDiscovered makes the check functions see a deterministic set of
// installed skills rather than whatever is on the developer's filesystem.
func stubDiscovered(t *testing.T, skills []installed.Skill) {
	t.Helper()
	old := discoverInstalled
	discoverInstalled = func() ([]installed.Skill, error) { return skills, nil }
	t.Cleanup(func() { discoverInstalled = old })
}

func bundledStale(name string) installed.Skill {
	stale := map[string][]byte{"SKILL.md": []byte("totally different")}
	return installed.Skill{
		Name:   name,
		Source: skill.SourceBundled,
		Files:  stale,
		Hash:   skill.ContentHash(stale),
	}
}

func TestBundledSkillUpdates_flagsDiverged(t *testing.T) {
	bs, err := bundled.All()
	require.NoError(t, err)
	require.NotEmpty(t, bs)

	stubDiscovered(t, []installed.Skill{bundledStale(bs[0].Name)})
	got := bundledSkillUpdates(newSkillCheckConfig(t, ""))
	assert.Equal(t, []string{bs[0].Name}, got)
}

func TestBundledSkillUpdates_ignoresUpToDate(t *testing.T) {
	bs, err := bundled.All()
	require.NoError(t, err)
	require.NotEmpty(t, bs)

	stubDiscovered(t, []installed.Skill{{
		Name:   bs[0].Name,
		Source: skill.SourceBundled,
		Files:  bs[0].Files,
		Hash:   skill.ContentHash(bs[0].Files),
	}})
	got := bundledSkillUpdates(newSkillCheckConfig(t, ""))
	assert.Empty(t, got)
}

func TestBundledSkillUpdates_ignoresRemoteSourceEntries(t *testing.T) {
	stubDiscovered(t, []installed.Skill{{
		Name:   "orbit",
		Source: skill.SourceRemote,
		Files:  map[string][]byte{"SKILL.md": []byte("x")},
		Hash:   skill.ContentHash(map[string][]byte{"SKILL.md": []byte("x")}),
	}})
	got := bundledSkillUpdates(newSkillCheckConfig(t, ""))
	assert.Empty(t, got, "remote-sourced skills should not appear in the bundled check")
}

func TestBundledSkillUpdates_configOptOut(t *testing.T) {
	bs, err := bundled.All()
	require.NoError(t, err)
	require.NotEmpty(t, bs)

	stubDiscovered(t, []installed.Skill{bundledStale(bs[0].Name)})
	got := bundledSkillUpdates(newSkillCheckConfig(t, "false"))
	assert.Empty(t, got)
}

func TestBundledSkillUpdates_envOptOut(t *testing.T) {
	t.Setenv("GLAB_NOTIFY_SKILL_UPDATES", "false")
	bs, err := bundled.All()
	require.NoError(t, err)
	require.NotEmpty(t, bs)

	stubDiscovered(t, []installed.Skill{bundledStale(bs[0].Name)})
	got := bundledSkillUpdates(newSkillCheckConfig(t, ""))
	assert.Empty(t, got)
}

func TestBundledSkillUpdates_deduplicatesAcrossScopes(t *testing.T) {
	bs, err := bundled.All()
	require.NoError(t, err)
	require.NotEmpty(t, bs)
	stale := map[string][]byte{"SKILL.md": []byte("stale")}

	stubDiscovered(t, []installed.Skill{
		{
			Name:   bs[0].Name,
			Source: skill.SourceBundled,
			Scope:  installed.ScopeProject,
			Files:  stale,
			Hash:   skill.ContentHash(stale),
		},
		{
			Name:   bs[0].Name,
			Source: skill.SourceBundled,
			Scope:  installed.ScopeGlobal,
			Files:  stale,
			Hash:   skill.ContentHash(stale),
		},
	})
	got := bundledSkillUpdates(newSkillCheckConfig(t, ""))
	assert.Equal(t, []string{bs[0].Name}, got, "same name in two scopes should only appear once")
}

func TestWriteSkillUpdateBlock_singularAndPlural(t *testing.T) {
	ios, _, _, stderr := cmdtest.TestIOStreams(cmdtest.WithTestIOStreamsAsTTY(false))

	writeSkillUpdateBlock(ios, nil)
	assert.Empty(t, stderr.String())

	stderr.Reset()
	writeSkillUpdateBlock(ios, []string{"glab"})
	got := stderr.String()
	assert.Contains(t, got, "Agent skill updates available")
	assert.Contains(t, got, "  glab\n")
	assert.Contains(t, got, "  Run: glab skills update glab\n")

	stderr.Reset()
	writeSkillUpdateBlock(ios, []string{"glab", "glab-stack"})
	got = stderr.String()
	assert.Contains(t, got, "Agent skill updates available")
	assert.Contains(t, got, "  glab, glab-stack\n")
	assert.Contains(t, got, "  Run: glab skills update --all\n")
}

func TestWriteSkillUpdateBlock_leadingBlankLineSeparatesFromPriorOutput(t *testing.T) {
	ios, _, _, stderr := cmdtest.TestIOStreams(cmdtest.WithTestIOStreamsAsTTY(false))
	writeSkillUpdateBlock(ios, []string{"glab"})
	assert.True(t, strings.HasPrefix(stderr.String(), "\n"),
		"block should begin with a blank line so it doesn't cram against preceding stderr output")
}
