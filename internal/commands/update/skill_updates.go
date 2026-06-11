package update

import (
	"strings"

	"gitlab.com/gitlab-org/cli/internal/commands/skills/bundled"
	"gitlab.com/gitlab-org/cli/internal/commands/skills/installed"
	"gitlab.com/gitlab-org/cli/internal/commands/skills/remote"
	"gitlab.com/gitlab-org/cli/internal/commands/skills/skill"
	"gitlab.com/gitlab-org/cli/internal/config"
	"gitlab.com/gitlab-org/cli/internal/dbg"
	"gitlab.com/gitlab-org/cli/internal/iostreams"
)

// Overridable so tests don't pick up the developer's ~/.agents/skills/.
var discoverInstalled = installed.Discover

// bundledSkillUpdates lists installed bundled skills whose on-disk content
// does not match the version embedded in this binary.
func bundledSkillUpdates(cfg config.Config) []string {
	return skillUpdates(cfg, skill.SourceBundled, bundled.Get)
}

// remoteSkillUpdates lists installed remote skills whose on-disk content
// does not match the current upstream. Each name triggers a gitlab.com
// request — gate on the 24h CheckUpdate cadence, not per command.
func remoteSkillUpdates(cfg config.Config) []string {
	return skillUpdates(cfg, skill.SourceRemote, remote.Get)
}

// skillUpdates is best-effort: discovery or getter failures return nil
// rather than an error so a stale check can't disrupt the user's command.
// Surfaced via dbg.Debugf so the silent path is at least visible under
// GLAB_DEBUG=1.
func skillUpdates(cfg config.Config, source skill.Source, getSource func(string) (skill.Skill, error)) []string {
	if !isSkillNotificationsEnabled(cfg) {
		return nil
	}
	all, err := discoverInstalled()
	if err != nil {
		dbg.Debugf("skill update check (%s): discover failed: %s", source, err)
		return nil
	}
	var out []string
	seen := map[string]struct{}{}
	for _, ins := range all {
		if ins.Source != source {
			continue
		}
		if _, ok := seen[ins.Name]; ok {
			continue
		}
		src, err := getSource(ins.Name)
		if err != nil {
			dbg.Debugf("skill update check (%s): fetch %q failed: %s", source, ins.Name, err)
			continue
		}
		if skill.ContentHash(src.Files) != ins.Hash {
			out = append(out, ins.Name)
			seen[ins.Name] = struct{}{}
		}
	}
	return out
}

// writeSkillUpdateBlock prints the agent-skill update nudge as its own
// block, matching the shape of writeHumanUpdateBlock and the post-upgrade
// banner: a colored header, an indented body listing the skills, and an
// indented `Run:` line. A leading blank line separates it from whatever
// preceded it on stderr.
func writeSkillUpdateBlock(io *iostreams.IOStreams, names []string) {
	if len(names) == 0 {
		return
	}
	c := io.Color()
	action := "glab skills update " + names[0]
	if len(names) > 1 {
		action = "glab skills update --all"
	}
	io.LogError("")
	io.LogError(c.Yellow("Agent skill updates available"))
	io.LogErrorf("  %s\n", strings.Join(names, ", "))
	io.LogErrorf("  Run: %s\n", action)
}

func isSkillNotificationsEnabled(cfg config.Config) bool {
	// cfg.Get already consults GLAB_NOTIFY_SKILL_UPDATES via EnvKeyEquivalence.
	val, err := cfg.Get("", "notify_skill_updates")
	if err != nil || val == "" {
		return true
	}
	switch strings.ToLower(strings.TrimSpace(val)) {
	case "false", "no", "n", "0":
		return false
	}
	return true
}
