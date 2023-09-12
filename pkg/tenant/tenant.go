package tenant

import (
	"fmt"
	"math/rand"
	"path"
	"strings"

	"github.com/clerkinc/clerk-sdk-go/clerk"
	"github.com/flanksource/tenant-controller/api/v1"
	"github.com/flanksource/tenant-controller/pkg/config"
	"github.com/flanksource/tenant-controller/pkg/utils"
	goslug "github.com/gosimple/slug"
)

func NewTenant(req v1.TenantRequestBody) (v1.Tenant, error) {
	// TODO: Hardcoded for now
	cloud := v1.AZURE

	kPath, err := utils.Template(config.Config.Git.KustomizationPath, map[string]any{
		"cluster": getClusterName(cloud),
	})
	if err != nil {
		return v1.Tenant{}, err
	}

	slug := req.Data.Slug
	if slug == "" {
		// TODO: If slug is empty, we might have to update the new slug in clerk
		slug = goslug.Make(req.Data.Name)
	}

	// For keeping postgres user and database name simple
	slug = strings.ReplaceAll(slug, "-", "_")

	return v1.Tenant{
		Name:              req.Data.Name,
		OrgID:             req.Data.OrgID,
		Cloud:             cloud,
		Slug:              slug,
		KustomizationPath: kPath,
		ContentPath:       path.Join(path.Dir(kPath), slug),
		Host:              getHost(cloud, slug),
		DBUsername:        fmt.Sprintf("%s_%d", strings.ToLower(slug), rand.Intn(1000)),
		DBPassword:        utils.RandomString(16),
	}, nil
}

func getClusterName(cloud v1.CloudProvider) string {
	// TODO: Take this from config
	switch cloud {
	case v1.AZURE:
		return "azure-internal-prod"
	case v1.AWS:
		return "aws-demo"
	}
	return ""
}

func getHost(cloud v1.CloudProvider, tenantName string) string {
	switch cloud {
	case v1.AZURE:
		return fmt.Sprintf("mission-control.%s.internal-prod.flanksource.com", tenantName)
	case v1.AWS:
		return ""
	default:
		return ""
	}
}

func updateHostOnClerk(orgID, host string) error {
	client, err := clerk.NewClient(config.Config.Clerk.SecretKey)
	if err != nil {
		return err
	}

	params := clerk.UpdateOrganizationMetadataParams{
		PublicMetadata: []byte(fmt.Sprintf(`{"backend_url": "https://%s"}`, host)),
	}
	if _, err := client.Organizations().UpdateMetadata(orgID, params); err != nil {
		return err
	}

	return nil
}
