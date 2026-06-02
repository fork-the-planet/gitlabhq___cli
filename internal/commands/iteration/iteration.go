package iteration

import (
	"github.com/MakeNowJust/heredoc/v2"
	"github.com/spf13/cobra"

	"gitlab.com/gitlab-org/cli/internal/cmdutils"
	iterationListCmd "gitlab.com/gitlab-org/cli/internal/commands/iteration/list"
)

func NewCmdIteration(f cmdutils.Factory) *cobra.Command {
	iterationCmd := &cobra.Command{
		Use:   "iteration <command> [flags]",
		Short: `Retrieve iteration information.`,
		Long: heredoc.Docf(`
			Iterations are time-boxed periods, similar to sprints, that group
			issues and merge requests in a project or group.

			List the iterations for the current project, or use the %[1]s--group%[1]s
			flag to list a group's iterations instead.
		`, "`"),
	}

	cmdutils.EnableRepoOverride(iterationCmd, f)

	iterationCmd.AddCommand(iterationListCmd.NewCmdList(f))
	return iterationCmd
}
