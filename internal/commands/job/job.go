package job

import (
	"github.com/MakeNowJust/heredoc/v2"
	"github.com/spf13/cobra"

	"gitlab.com/gitlab-org/cli/internal/cmdutils"
	"gitlab.com/gitlab-org/cli/internal/commands/job/artifact"
)

func NewCmdJob(f cmdutils.Factory) *cobra.Command {
	jobCmd := &cobra.Command{
		Use:   "job <command> [flags]",
		Short: `Work with GitLab CI/CD jobs.`,
		Long: heredoc.Docf(`
			Inspect CI/CD jobs from a pipeline. Download the artifacts a job
			produced with %[1]sjob artifact%[1]s.

			Use %[1]s--repo%[1]s to target a project other than the current one.
		`, "`"),
	}

	cmdutils.EnableRepoOverride(jobCmd, f)
	jobCmd.AddCommand(artifact.NewCmdArtifact(f))
	return jobCmd
}
