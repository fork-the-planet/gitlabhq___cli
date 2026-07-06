package security

import (
	"github.com/MakeNowJust/heredoc/v2"
	"github.com/spf13/cobra"

	"gitlab.com/gitlab-org/cli/internal/cmdutils"
	configCmd "gitlab.com/gitlab-org/cli/internal/commands/security/config"
	"gitlab.com/gitlab-org/cli/internal/text"
)

func NewCmd(f cmdutils.Factory) *cobra.Command {
	securityCmd := &cobra.Command{
		Use:   "security <command> [flags]",
		Short: "Manage GitLab security features. (EXPERIMENTAL)",
		Long: heredoc.Doc(`
			Configure GitLab security features for a project.
		`) + text.ExperimentalString,
	}

	cmdutils.EnableRepoOverride(securityCmd, f)

	securityCmd.AddCommand(configCmd.NewCmd(f))

	return securityCmd
}
