package opentofu

import (
	"github.com/MakeNowJust/heredoc/v2"
	"github.com/spf13/cobra"

	"gitlab.com/gitlab-org/cli/internal/cmdutils"
	initCmd "gitlab.com/gitlab-org/cli/internal/commands/opentofu/init"
	stateCmd "gitlab.com/gitlab-org/cli/internal/commands/opentofu/state"
)

func NewCmd(f cmdutils.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "opentofu <command> [flags]",
		Short: `Work with the OpenTofu or Terraform integration.`,
		Long: heredoc.Doc(`
			Manage OpenTofu and Terraform state stored in the GitLab-managed
			state backend.

			Initialize a working directory against a GitLab-managed state, and
			list, lock, unlock, download, or delete states.
		`),
		Aliases: []string{"terraform", "tf"},
	}

	cmdutils.EnableRepoOverride(cmd, f)

	cmd.AddCommand(initCmd.NewCmd(f))
	cmd.AddCommand(stateCmd.NewCmd(f))
	return cmd
}
