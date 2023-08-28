package connectors

import (
	"context"
	"strings"

	v1 "github.com/flanksource/tenant-controller/api/v1"
	"github.com/go-git/go-billy/v5"
	git "github.com/go-git/go-git/v5"
	"github.com/pkg/errors"
)

type Connector interface {
	Clone(ctx context.Context, branch, local string) (billy.Filesystem, *git.Worktree, error)
	Push(ctx context.Context, branch string) error
	OpenPullRequest(ctx context.Context, base string, head string, spec *v1.PullRequestTemplate) (int, error)
	ClosePullRequest(ctx context.Context, id int) error
}

func NewConnector(git_config *v1.GitopsAPISpec) (Connector, error) {
	if strings.HasPrefix(git_config.Repository, "https://github.com/") {
		path := git_config.Repository[19:]
		parts := strings.Split(path, "/")
		if len(parts) != 2 {
			return nil, errors.Errorf("invalid repository url: %s", git_config.Repository)
		}
		owner := parts[0]
		repoName := parts[1]
		repoName = strings.TrimSuffix(repoName, ".git")
		githubToken := git_config.GITHUB_TOKEN
		return NewGithub(owner, repoName, githubToken)
	} else if strings.HasPrefix(git_config.Repository, "ssh://") {
		sshURL := git_config.Repository[6:]
		user := strings.Split(sshURL, "@")[0]

		privateKey := git_config.SSH_PRIVATE_KEY
		password := git_config.SSH_PRIVATE_KEY_PASSORD
		return NewGitSSH(sshURL, user, []byte(privateKey), password)
	}
	return nil, errors.New("no connector settings found")
}
