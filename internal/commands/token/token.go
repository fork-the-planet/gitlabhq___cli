package token

import (
	"github.com/MakeNowJust/heredoc/v2"
	"github.com/spf13/cobra"

	"gitlab.com/gitlab-org/cli/internal/cmdutils"
	"gitlab.com/gitlab-org/cli/internal/commands/token/create"
	"gitlab.com/gitlab-org/cli/internal/commands/token/list"
	"gitlab.com/gitlab-org/cli/internal/commands/token/revoke"
	"gitlab.com/gitlab-org/cli/internal/commands/token/rotate"
)

func NewTokenCmd(f cmdutils.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "token",
		Short: "Manage personal, project, or group tokens.",
		Long: heredoc.Doc(`
			Create, list, revoke, and rotate access tokens for your user
			account, a project, or a group.
		`),
		Aliases: []string{"token"},
	}

	cmdutils.EnableRepoOverride(cmd, f)
	cmd.AddCommand(create.NewCmdCreate(f))
	cmd.AddCommand(revoke.NewCmdRevoke(f))
	cmd.AddCommand(rotate.NewCmdRotate(f))
	cmd.AddCommand(list.NewCmdList(f))
	return cmd
}
