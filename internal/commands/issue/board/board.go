package board

import (
	"github.com/MakeNowJust/heredoc/v2"
	"github.com/spf13/cobra"

	"gitlab.com/gitlab-org/cli/internal/cmdutils"
	boardCreateCmd "gitlab.com/gitlab-org/cli/internal/commands/issue/board/create"
	boardViewCmd "gitlab.com/gitlab-org/cli/internal/commands/issue/board/view"
)

func NewCmdBoard(f cmdutils.Factory) *cobra.Command {
	issueCmd := &cobra.Command{
		Use:   "board [command] [flags]",
		Short: `Work with GitLab issue boards in the given project.`,
		Long: heredoc.Doc(`
			Issue boards organize issues into lists, similar to a kanban board,
			so you can track work as it moves between stages. Create a board, or
			view one interactively in the terminal.
		`),
	}

	issueCmd.AddCommand(boardCreateCmd.NewCmdCreate(f))
	issueCmd.AddCommand(boardViewCmd.NewCmdView(f))
	issueCmd.PersistentFlags().StringP("repo", "R", "", "Select another repository using the OWNER/REPO format or the project ID. Supports group namespaces.")

	return issueCmd
}
