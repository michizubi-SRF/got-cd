package helper

import (
	"fmt"
	"net/url"
	"os/exec"
	"strings"
)

// GetCurrentOwnerAndRepo returns the GitHub owner (user or organization) and
// repository name of the current repository, derived from the "origin" remote.
func GetCurrentOwnerAndRepo() (string, string, error) {
	output, err := exec.Command("git", "remote", "get-url", "origin").Output()
	if err != nil {
		return "", "", fmt.Errorf("error getting remote URL of origin: %w", err)
	}
	return ParseOwnerAndRepo(string(output))
}

// ParseOwnerAndRepo extracts owner and repository name from a git remote URL.
// Supported formats:
//
//	git@github.com:owner/repo.git
//	https://github.com/owner/repo(.git)
//	ssh://git@github.com(:port)/owner/repo(.git)
func ParseOwnerAndRepo(remoteURL string) (string, string, error) {
	remoteURL = strings.TrimSpace(remoteURL)

	var path string
	if strings.Contains(remoteURL, "://") {
		u, err := url.Parse(remoteURL)
		if err != nil {
			return "", "", fmt.Errorf("invalid remote URL %q: %w", remoteURL, err)
		}
		path = u.Path
	} else if at := strings.Index(remoteURL, "@"); at >= 0 && strings.Contains(remoteURL[at:], ":") {
		// scp-like syntax: user@host:owner/repo.git
		path = remoteURL[strings.Index(remoteURL[at:], ":")+at+1:]
	} else {
		return "", "", fmt.Errorf("unsupported remote URL %q", remoteURL)
	}

	path = strings.Trim(path, "/")
	path = strings.TrimSuffix(path, ".git")

	parts := strings.Split(path, "/")
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		return "", "", fmt.Errorf("could not determine owner and repository from remote URL %q", remoteURL)
	}

	return parts[0], parts[1], nil
}
