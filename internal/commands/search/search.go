package search

import (
	"github.com/MakeNowJust/heredoc/v2"
	"github.com/spf13/cobra"

	"gitlab.com/gitlab-org/cli/internal/cmdutils"
	semanticCmd "gitlab.com/gitlab-org/cli/internal/commands/search/semantic"
	"gitlab.com/gitlab-org/cli/internal/text"
)

func NewCmd(f cmdutils.Factory) *cobra.Command {
	searchCmd := &cobra.Command{
		Use:   "search <command> [flags]",
		Short: `Search for code and resources in a GitLab project. (BETA)`,
		Long: heredoc.Docf(`
			Search a GitLab project for code and other resources. The
			%[1]ssemantic%[1]s subcommand runs an AI-powered semantic code search.

			Use %[1]s--repo%[1]s to target a project other than the current one.
		`, "`") + text.BetaString,
	}

	cmdutils.EnableRepoOverride(searchCmd, f)

	searchCmd.AddCommand(semanticCmd.NewCmd(f))

	return searchCmd
}
