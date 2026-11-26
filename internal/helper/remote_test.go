package helper

import "testing"

func TestParseOwnerAndRepo(t *testing.T) {
	tests := []struct {
		url   string
		owner string
		repo  string
	}{
		{"git@github.com:michizubi-SRF/got-cd.git", "michizubi-SRF", "got-cd"},
		{"git@github.com:other-org/some.repo.git", "other-org", "some.repo"},
		{"git@github.com:org/repo", "org", "repo"},
		{"https://github.com/org/repo.git", "org", "repo"},
		{"https://github.com/org/repo", "org", "repo"},
		{"https://github.com/org/repo/", "org", "repo"},
		{"https://user@github.com/org/repo.git", "org", "repo"},
		{"ssh://git@github.com/org/repo.git", "org", "repo"},
		{"ssh://git@github.com:22/org/repo.git", "org", "repo"},
		{"  git@github.com:org/repo.git\n", "org", "repo"},
	}

	for _, tt := range tests {
		owner, repo, err := ParseOwnerAndRepo(tt.url)
		if err != nil {
			t.Errorf("ParseOwnerAndRepo(%q) returned error: %v", tt.url, err)
			continue
		}
		if owner != tt.owner || repo != tt.repo {
			t.Errorf("ParseOwnerAndRepo(%q) = (%q, %q), want (%q, %q)", tt.url, owner, repo, tt.owner, tt.repo)
		}
	}
}

func TestParseOwnerAndRepoInvalid(t *testing.T) {
	for _, url := range []string{"", "not a url", "https://github.com/onlyowner", "git@github.com:repo.git"} {
		if _, _, err := ParseOwnerAndRepo(url); err == nil {
			t.Errorf("ParseOwnerAndRepo(%q) expected error, got nil", url)
		}
	}
}
