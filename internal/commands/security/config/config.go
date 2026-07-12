package config

import (
	"github.com/MakeNowJust/heredoc/v2"
	"github.com/spf13/cobra"

	"gitlab.com/gitlab-org/cli/internal/cmdutils"
	disableCmd "gitlab.com/gitlab-org/cli/internal/commands/security/config/disable"
	enableCmd "gitlab.com/gitlab-org/cli/internal/commands/security/config/enable"
	statusCmd "gitlab.com/gitlab-org/cli/internal/commands/security/config/status"
	"gitlab.com/gitlab-org/cli/internal/text"
)

func NewCmd(f cmdutils.Factory) *cobra.Command {
	configCmd := &cobra.Command{
		Use:   "config <command> [flags]",
		Short: "Configure security scan profiles for a project. (EXPERIMENTAL)",
		Long: heredoc.Doc(`
			Enable, disable, or inspect security scan profiles for a project.

			A profile bundles a set of security scans, such as SAST, secret
			detection, dependency scanning, or container scanning, or post-scan
			processing on given scans, like dependency scanning auto remediation.
		`) + text.ExperimentalString,
	}

	configCmd.AddCommand(enableCmd.NewCmd(f))
	configCmd.AddCommand(disableCmd.NewCmd(f))
	configCmd.AddCommand(statusCmd.NewCmd(f))

	return configCmd
}
