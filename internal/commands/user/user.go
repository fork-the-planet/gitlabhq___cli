package user

import (
	"github.com/MakeNowJust/heredoc/v2"
	"github.com/spf13/cobra"

	"gitlab.com/gitlab-org/cli/internal/cmdutils"
	userEventsCmd "gitlab.com/gitlab-org/cli/internal/commands/user/events"
)

func NewCmdUser(f cmdutils.Factory) *cobra.Command {
	userCmd := &cobra.Command{
		Use:   "user <command> [flags]",
		Short: "Interact with a GitLab user account.",
		Long: heredoc.Doc(`
			Look up information about GitLab users, such as the events a user has
			generated.
		`),
	}

	userCmd.AddCommand(userEventsCmd.NewCmdEvents(f))

	return userCmd
}
