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

func CreateTenantResources(connector connectors.Connector, tenant *pkg.Tenant, tenantObjs []*unstructured.Unstructured) (work *gitv5.Worktree, title string, err error) {
	title = fmt.Sprintf("feat: add %s tenant resources", tenant.Name)
	pkg.Config.GIT.SetDefaults(title)
	fs, work, err := connector.Clone(context.TODO(), pkg.Config.GIT.Base, pkg.Config.GIT.Branch)
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
		if err = copy(body, contentPath, fs, work); err != nil {
			return nil, "", err
		}
	}
	// update root kustomization and add tenant kustomization to it
	kustomization, err := getKustomizaton(fs, pkg.Config.GIT.Kustomization)
	if err != nil {
		return nil, "", err
	}
	kustomization.Resources = append(kustomization.Resources, tenant.Slug)
	existingKustomization, err := yaml.Marshal(kustomization)
	if err != nil {
		return nil, "", err
	}
	if err = copy(existingKustomization, pkg.Config.GIT.Kustomization, fs, work); err != nil {
		return nil, "", err
	}
	return
}

func CreateCommit(work *gitv5.Worktree, title string) (hash string, err error) {
	author := &object.Signature{
		Name:  pkg.Config.GIT.User,
		Email: pkg.Config.GIT.Email,
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
	pkg.Config.GIT.Kustomization, _ = pkg.Template(pkg.Config.GIT.Kustomization, map[string]interface{}{
		"cluster": getClusterName(tenant),
	})
	return path.Dir(pkg.Config.GIT.Kustomization) + "/" + tenant.Slug
}

func getClusterName(tenant *pkg.Tenant) string {
	switch tenant.Cloud {
	case pkg.AZURE:
		return "azure-internal-prod"
	case pkg.AWS:
		return "aws-demo"
	}
	return ""
}

func copy(data []byte, path string, fs billy.Filesystem, work *gitv5.Worktree) error {
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
