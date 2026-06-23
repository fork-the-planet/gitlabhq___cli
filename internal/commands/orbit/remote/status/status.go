package status

import (
	"context"

	"github.com/MakeNowJust/heredoc/v2"
	"github.com/spf13/cobra"

	gitlab "gitlab.com/gitlab-org/api/client-go/v2"

	"gitlab.com/gitlab-org/cli/internal/api"
	"gitlab.com/gitlab-org/cli/internal/cmdutils"
	"gitlab.com/gitlab-org/cli/internal/commands/orbit/internal/orbiterr"
	"gitlab.com/gitlab-org/cli/internal/iostreams"
	"gitlab.com/gitlab-org/cli/internal/mcpannotations"
	"gitlab.com/gitlab-org/cli/internal/text"
)

type options struct {
	apiClient func(repoHost string) (*api.Client, error)
	io        *iostreams.IOStreams

	hostname string
}

func NewCmd(f cmdutils.Factory) *cobra.Command {
	opts := &options{
		apiClient: f.ApiClient,
		io:        f.IO(),
	}

	cmd := &cobra.Command{
		Use:   "status",
		Short: `Show GitLab Knowledge Graph cluster health. (EXPERIMENTAL)`,
		Long: heredoc.Doc(`
			Prints the Orbit cluster health as pretty-printed JSON. Use this
			command to confirm Orbit is enabled and reachable for your user.
			It is the first step in the Orbit discovery workflow.

			The output is always the system health object (fields such as
			status, version, components, error) regardless of the GitLab
			instance version. On newer instances (19.1+), when the backend
			cannot reach the gRPC cluster, the output includes an error
			field (and status is "unknown"). Use --jq to filter health
			fields directly, for example --jq '.status' or
			--jq '.components'. The user/system wrapper returned by newer
			instances is not exposed in the output.
		`) + text.ExperimentalString,
		Example: heredoc.Doc(`
			$ glab orbit remote status
		`),
		Annotations: map[string]string{
			mcpannotations.Safe: "true",
		},
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return opts.run(cmd.Context())
		},
	}

	cmd.Flags().StringVar(&opts.hostname, "hostname", "",
		"GitLab hostname to query. Defaults to the current repository's host or `gitlab.com`.")

	cmdutils.AddJQFlag(cmd, f.IO())
	return cmd
}

func (o *options) run(ctx context.Context) error {
	client, err := o.apiClient(o.hostname)
	if err != nil {
		return err
	}

	status, _, err := client.Lab().Orbit.GetStatus(nil, gitlab.WithContext(ctx))
	if err != nil {
		return orbiterr.Translate(err)
	}

	// New nested shape: the API returned a "user" wrapper.
	if status.User != nil {
		if !status.User.Available {
			return orbiterr.UnavailableForUser()
		}
		if status.System != nil {
			return o.io.PrintJSON(status.System)
		}
		// User has access but the server returned no system health
		// object. Under the current API contract this is unreachable
		// (System is present whenever Available is true), but we
		// surface an explicit error rather than silently printing
		// the user/system wrapper.
		return orbiterr.SystemHealthAbsent()
	}

	// Old flat shape (pre-19.1 instances): print the whole
	// status object as-is.
	return o.io.PrintJSON(status)
}
