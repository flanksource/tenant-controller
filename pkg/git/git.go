package git

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/flanksource/tenant-controller/pkg"
	"github.com/flanksource/tenant-controller/pkg/git/connectors"
	"github.com/go-git/go-billy/v5"
	gitv5 "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"sigs.k8s.io/kustomize/api/types"
)

func OpenPRWithTenantResources(tenant *pkg.Tenant, tenantObjs []*unstructured.Unstructured) (pr int, hash string, err error) {
	// TODO: Git config should be passed as arg, not taken from global
	connector, err := connectors.NewConnector(pkg.Config.Git)
	if err != nil {
		return
	}

	title := fmt.Sprintf("feat: add %s tenant resources", tenant.Name)
	work, title, err := CreateTenantResources(connector, tenant, tenantObjs)
	if err != nil {
		return
	}

	hash, err = CreateCommit(work, title)
	if err != nil {
		return
	}

	if err = connector.Push(context.Background(), fmt.Sprintf("%s:%s", pkg.Config.Git.Branch, pkg.Config.Git.Base)); err != nil {
		return
	}

	pr, err = connector.OpenPullRequest(context.Background(), pkg.Config.Git.Base, pkg.Config.Git.Branch, pkg.Config.Git.PullRequest)
	if err != nil {
		return
	}

	return
}

func CreateTenantResources(connector connectors.Connector, tenant *pkg.Tenant, tenantObjs []*unstructured.Unstructured) (work *gitv5.Worktree, title string, err error) {
	// TODO: This is not thread safe
	pkg.Config.Git.SetDefaults(title)

	fs, work, err := connector.Clone(context.Background(), pkg.Config.Git.Base, pkg.Config.Git.Branch)
	if err != nil {
		return
	}
	// path at which tenant resources will be stored
	contentsPath := getContentsPath(tenant)

	// add tenant resources to git
	for _, obj := range tenantObjs {
		contentPath := filepath.Join(contentsPath, strings.ToLower(obj.GetKind())+".yaml")
		body, err := yaml.Marshal(obj.Object)
		if err != nil {
			return nil, "", err
		}
		if err = writeGitWorkTree(body, contentPath, fs, work); err != nil {
			return nil, "", err
		}
	}
	// update root kustomization and add tenant kustomization to it
	kustomization, err := getKustomizaton(fs, pkg.Config.Git.KustomizationPath)
	if err != nil {
		return nil, "", err
	}

	// TODO: This should not append the resources, tenant yaml files should be in
	// their own directories
	kustomization.Resources = append(kustomization.Resources, tenant.Slug)
	existingKustomization, err := yaml.Marshal(kustomization)
	if err != nil {
		return nil, "", err
	}
	if err = writeGitWorkTree(existingKustomization, pkg.Config.Git.KustomizationPath, fs, work); err != nil {
		return nil, "", err
	}
	return
}

func CreateCommit(work *gitv5.Worktree, title string) (hash string, err error) {
	author := &object.Signature{
		Name:  pkg.Config.Git.User,
		Email: pkg.Config.Git.Email,
		When:  time.Now(),
	}
	if author.Name == "" {
		author.Name = "Tenant Operator"
	}
	if author.Email == "" {
		author.Email = "noreply@tenant-operator"
	}
	_hash, err := work.Commit(title, &gitv5.CommitOptions{
		Author: author,
		All:    true,
	})

	if err != nil {
		return
	}
	hash = _hash.String()
	return
}

func getContentsPath(tenant *pkg.Tenant) string {
	pkg.Config.Git.KustomizationPath, _ = pkg.Template(pkg.Config.Git.KustomizationPath, map[string]interface{}{
		"cluster": getClusterName(tenant),
	})
	return path.Dir(pkg.Config.Git.KustomizationPath) + "/" + tenant.Slug
}

func getClusterName(tenant *pkg.Tenant) string {
	switch tenant.Cloud {
	case pkg.Azure:
		return "azure-internal-prod"
	case pkg.AWS:
		return "aws-demo"
	}
	return ""
}

func writeGitWorkTree(data []byte, path string, fs billy.Filesystem, work *gitv5.Worktree) error {
	dst, err := openOrCreate(path, fs)
	if err != nil {
		return errors.Wrap(err, "failed to open")
	}
	src := bytes.NewBuffer(data)
	if _, err := io.Copy(dst, src); err != nil {
		return errors.Wrap(err, "failed to copy")
	}
	if err := dst.Close(); err != nil {
		return errors.Wrap(err, "failed to close")
	}
	_, err = work.Add(path)
	return errors.Wrap(err, "failed to add to git")
}

func openOrCreate(path string, fs billy.Filesystem) (billy.File, error) {
	return fs.Create(path)
}

func getKustomizaton(fs billy.Filesystem, path string) (*types.Kustomization, error) {
	kustomization := types.Kustomization{}

	if _, err := fs.Stat(path); err != nil {
		return &kustomization, nil
	}
	existing, err := fs.Open(path)
	if err != nil {
		return nil, err
	}
	existingKustomization, err := io.ReadAll(existing)
	if err != nil {
		return nil, err
	}
	if err := yaml.Unmarshal(existingKustomization, &kustomization); err != nil {
		return nil, err
	}
	return &kustomization, nil
}
