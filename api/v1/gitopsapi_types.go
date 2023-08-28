package v1

// GitopsAPISpec defines the desired state of GitopsAPI
type GitopsAPISpec struct {
	// The repository URL, can be a HTTP or SSH address.
	Repository string `json:"repository,omitempty"`
	User       string `json:"user,omitempty"`
	Email      string `json:"email,omitempty"`

	// Open a new Pull request from the branch back to the base
	PullRequest *PullRequestTemplate `json:"pull_request,omitempty"`

	// For Github repositories it must contain GITHUB_TOKEN
	// +optional
	GITHUB_TOKEN string `json:"github_token,omitempty"`

	// For SSH repositories the secret must contain SSH_PRIVATE_KEY, SSH_PRIVATE_KEY_PASSORD
	SSH_PRIVATE_KEY         string `json:"ssh_private_key,omitempty"`
	SSH_PRIVATE_KEY_PASSORD string `json:"ssh_private_key_password,omitempty"`

	// The path to a kustomization file to insert or remove the resource, can included templated values .e.g `specs/clusters/{{.cluster}}/kustomization.yaml`
	KustomizationPath string `json:"kustomizationPath,omitempty"`
}

type PullRequestTemplate struct {
	Body      string   `json:"body,omitempty"`
	Title     string   `json:"title,omitempty"`
	Reviewers []string `json:"reviewers,omitempty"`
	Assignees []string `json:"assignees,omitempty"`
	Tags      []string `json:"tags,omitempty"`

	// The branch to use as a baseline for the new branch, defaults to master
	Base string `json:"base,omitempty"`
	// The branch to push updates back to, defaults to master
	Branch string `json:"branch,omitempty"`
}
