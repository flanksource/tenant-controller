package tenant

import (
	"encoding/json"
	"fmt"
	"path"
	"strings"

	"github.com/clerkinc/clerk-sdk-go/clerk"
	"github.com/flanksource/tenant-controller/api/v1"
	"github.com/flanksource/tenant-controller/pkg/utils"
	goslug "github.com/gosimple/slug"
)

const (
	TenantStateActive    = "active"
	TenantStateSuspended = "suspended"
)

func NewTenant(req v1.TenantRequestBody) (v1.Tenant, error) {
	kPath, err := utils.Template(v1.GlobalConfig.Git.KustomizationPath, map[string]any{
		"cluster": v1.GlobalConfig.GetClusterName(),
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
		Cloud:             v1.GlobalConfig.DefaultCloud,
		Slug:              slug,
		KustomizationPath: kPath,
		ContentPath:       path.Join(path.Dir(kPath), id),
		Host:              v1.GlobalConfig.GetHost(id),
		DBUsername:        strings.ToLower(orgID),
		DBPassword:        utils.RandomString(16),
	}, nil
}

func updateParamsOnClerk(tenant v1.Tenant) error {
	client, err := clerk.NewClient(v1.GlobalConfig.Clerk.SecretKey)
	if err != nil {
		return fmt.Errorf("error creating clerk client: %w", err)
	}

	pubMetadata, err := json.Marshal(map[string]string{
		"backend_url": fmt.Sprintf("https://%s", tenant.Host),
		"tenant_id":   tenant.ID,
		"state":       TenantStateActive,
	})
	if err != nil {
		return fmt.Errorf("error marshaling public metadata to json: %w", err)
	}
	if _, err := client.Organizations().Update(tenant.OrgID, clerk.UpdateOrganizationParams{
		Slug:           &tenant.Slug,
		PublicMetadata: pubMetadata,
	}); err != nil {
		return fmt.Errorf("error updating org on clerk: %w", err)
	}

	return nil
}
