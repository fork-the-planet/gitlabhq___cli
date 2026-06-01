package ssh

import (
	"github.com/MakeNowJust/heredoc/v2"
	"github.com/spf13/cobra"

	"gitlab.com/gitlab-org/cli/internal/cmdutils"
	cmdAdd "gitlab.com/gitlab-org/cli/internal/commands/ssh-key/add"
	cmdDelete "gitlab.com/gitlab-org/cli/internal/commands/ssh-key/delete"
	cmdGet "gitlab.com/gitlab-org/cli/internal/commands/ssh-key/get"
	cmdList "gitlab.com/gitlab-org/cli/internal/commands/ssh-key/list"
)

func NewCmdSSHKey(f cmdutils.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "ssh-key <command>",
		Short: "Manage SSH keys registered with your GitLab account.",
		Long: heredoc.Doc(`
			Add, list, get, and delete the SSH keys associated with your account.

			GitLab uses SSH keys to authenticate Git operations over SSH, and,
			depending on each key's usage type, to verify signed commits.
		`),
	}

	cmdutils.EnableRepoOverride(cmd, f)

	cmd.AddCommand(cmdAdd.NewCmdAdd(f))
	cmd.AddCommand(cmdGet.NewCmdGet(f))
	cmd.AddCommand(cmdList.NewCmdList(f))
	cmd.AddCommand(cmdDelete.NewCmdDelete(f))

	return cmd
}
