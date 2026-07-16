package config

import (
	"errors"
	"fmt"
	"net"
	"net/url"
	"strings"
)

// Settings is the registry connection information every package manager
// needs to route traffic through the firewall. It is populated by
// RegistrySettings and consumed by both the manager-specific config
// writers (npmrc today; later stack slices add yarnrc, pip, maven, …)
// and by the pm.Run wrapper.
type Settings struct {
	// RegistryURL is the full https URL package managers should resolve
	// against, with a trailing slash.
	RegistryURL string
	// AuthHost is the "host + path" portion of RegistryURL, used by npm-
	// style config files that scope credentials by URL prefix (for
	// example, //gitlab.example/api/v4/projects/42/packages/npm/:_authToken).
	AuthHost string
	// AuthToken is the GitLab token to send to the registry, or empty
	// when the registry host does not match the logged-in GitLab host.
	// The empty case is deliberate: it keeps the token out of requests
	// to third-party or public registries.
	AuthToken string
}

// RegistrySettings builds the [Settings] for resolveURL, a full registry
// URL. The GitLab token is attached only when the registry host matches
// the logged-in GitLab host, so the token is never sent to a third-party
// (for example, public) registry. The returned settings are shared by all
// package managers.
func RegistrySettings(gitlabHost, resolveURL, token string) (Settings, error) {
	if resolveURL == "" {
		return Settings{}, errors.New("no resolve registry configured; run 'glab df <manager>-config --repo-resolve <registry-url>'")
	}

	u, err := url.Parse(resolveURL)
	if err != nil {
		return Settings{}, fmt.Errorf("invalid resolve registry URL %q: %w", resolveURL, err)
	}
	if !u.IsAbs() || (u.Scheme != "http" && u.Scheme != "https") || u.Host == "" {
		return Settings{}, fmt.Errorf("resolve registry must be an absolute http(s) URL, got %q", resolveURL)
	}

	if !strings.HasSuffix(u.Path, "/") {
		u.Path += "/"
		if u.RawPath != "" {
			u.RawPath += "/"
		}
	}

	authToken := ""
	if hostsMatch(u.Host, gitlabHost, u.Scheme) {
		authToken = token
	}

	return Settings{
		RegistryURL: u.String(),
		AuthHost:    u.Host + u.EscapedPath(),
		AuthToken:   authToken,
	}, nil
}

// hostsMatch reports whether two hosts refer to the same registry. The
// comparison is case-insensitive and treats the scheme's default port
// (443 for https, 80 for http) as equivalent to no port.
func hostsMatch(a, b, scheme string) bool {
	return normalizeHost(a, scheme) == normalizeHost(b, scheme)
}

// normalizeHost lowercases host and strips a trailing default port for the
// given scheme. Hosts without a port are returned lowercased unchanged.
func normalizeHost(host, scheme string) string {
	host = strings.ToLower(host)

	hostname, port, err := net.SplitHostPort(host)
	if err != nil {
		return host
	}

	defaultPort := ""
	switch scheme {
	case "https":
		defaultPort = "443"
	case "http":
		defaultPort = "80"
	}

	if port == defaultPort {
		return hostname
	}

	return host
}
