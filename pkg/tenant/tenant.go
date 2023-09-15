package tenant

import (
	"fmt"
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

	// Kubernetes namespaces cannot have `_`
	id := strings.Replace(req.Data.OrgID, "org_", "org-", 1)

	return v1.Tenant{
		Name:              req.Data.Name,
		OrgID:             id,
		Cloud:             cloud,
		Slug:              slug,
		KustomizationPath: kPath,
		ContentPath:       path.Join(path.Dir(kPath), id),
		Host:              getHost(cloud, id),
		DBUsername:        id,
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

func getHost(cloud v1.CloudProvider, tenantID string) string {
	switch cloud {
	case v1.AZURE:
		return fmt.Sprintf("mission-control.%s.internal-prod.flanksource.com", tenantID)
	case v1.AWS:
		return ""
	default:
		return ""
	}
}

func updateParamsOnClerk(tenant v1.Tenant) error {
	client, err := clerk.NewClient(config.Config.Clerk.SecretKey)
	if err != nil {
		return err
	}

	if _, err := client.Organizations().Update(tenant.OrgID, clerk.UpdateOrganizationParams{
		Slug:           &tenant.Slug,
		PublicMetadata: []byte(fmt.Sprintf(`{"backend_url": "https://%s"}`, tenant.Host)),
	}); err != nil {
		return err
	}

	return nil
}
