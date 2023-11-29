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

	orgID := req.Data.OrgID

	// Kubernetes namespaces cannot have `_`
	// We are taking the last 12 chars to due to constraints
	// on domain name length and kubenretes resource name length
	_id := strings.ToLower(strings.Replace(orgID, "org_", "", 1))
	id := fmt.Sprintf("org-%s", _id[len(_id)-12:])

	return v1.Tenant{
		ID:                id,
		Name:              req.Data.Name,
		OrgID:             orgID,
		Cloud:             cloud,
		Slug:              slug,
		KustomizationPath: kPath,
		ContentPath:       path.Join(path.Dir(kPath), id),
		Host:              getHost(cloud, id),
		DBUsername:        strings.ToLower(orgID),
		DBPassword:        utils.RandomString(16),
	}, nil
}

func getClusterName(cloud v1.CloudProvider) string {
	// TODO: Take this from config
	switch cloud {
	case v1.AZURE:
		return config.Config.Azure.TenantCluster
	case v1.AWS:
		return config.Config.AWS.TenantCluster
	}
	return ""
}

func getHost(cloud v1.CloudProvider, tenantID string) string {
	switch cloud {
	case v1.AZURE:
		return fmt.Sprintf(config.Config.Azure.TenantHostFormat, tenantID)
	case v1.AWS:
		return fmt.Sprintf(config.Config.AWS.TenantHostFormat, tenantID)
	default:
		return ""
	}
}

func updateParamsOnClerk(tenant v1.Tenant) error {
	client, err := clerk.NewClient(config.Config.Clerk.SecretKey)
	if err != nil {
		return fmt.Errorf("error creating clerk client: %w", err)
	}

	if _, err := client.Organizations().Update(tenant.OrgID, clerk.UpdateOrganizationParams{
		Slug:           &tenant.Slug,
		PublicMetadata: []byte(fmt.Sprintf(`{"backend_url": "https://%s"}`, tenant.Host)),
	}); err != nil {
		return fmt.Errorf("error updating org on clerk: %w", err)
	}

	return nil
}
