package todo

import (
	"github.com/MakeNowJust/heredoc/v2"
	"github.com/spf13/cobra"

	"gitlab.com/gitlab-org/cli/internal/cmdutils"
	doneCmd "gitlab.com/gitlab-org/cli/internal/commands/todo/done"
	listCmd "gitlab.com/gitlab-org/cli/internal/commands/todo/list"
)

func NewCmd(f cmdutils.Factory) *cobra.Command {
	todoCmd := &cobra.Command{
		Use:   "todo <command> [flags]",
		Short: "Manage your to-do list.",
		Long: heredoc.Doc(`
			Your to-do list collects the items that need your attention, such as
			issues and merge requests where you were assigned, mentioned, or
			asked to review.

			List pending items, then mark them as done individually or all at once.
		`),
		Example: heredoc.Doc(`
			glab todo list
			glab todo done 123
			glab todo done --all
		`),
	}

	todoCmd.AddCommand(listCmd.NewCmd(f))
	todoCmd.AddCommand(doneCmd.NewCmd(f))

	return todoCmd
}
