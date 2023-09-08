package pkg

import (
	"fmt"
	"math/rand"
	"path"
	"strings"

	goslug "github.com/gosimple/slug"
)

type CloudProvider string

const (
	AWS   CloudProvider = "aws"
	Azure CloudProvider = "azure"
)

// Sample organization.created event payload from Clerk
//  {
//    "data": {
//      "created_at": 1654013202977,
//      "created_by": "user_1vq84bqWzw7qmFgqSwN4CH1Wp0n",
//      "id": "org_29w9IfBrPmcpi0IeBVaKtA7R94W",
//      "image_url": "https://img.clerk.com/xxxxxx",
//      "logo_url": "https://example.org/example.png",
//      "name": "Acme Inc",
//      "object": "organization",
//      "public_metadata": {},
//      "slug": "acme-inc",
//      "updated_at": 1654013202977
//    },
//    "object": "event",
//    "type": "organization.created"
//  }

type TenantRequestBody struct {
	Type   string `json:"type"`
	Object string `json:"object"`
	Data   struct {
		Slug  string `json:"slug"`
		OrgID string `json:"id"`
		Name  string `json:"name"`
	} `json:"data"`
}

type Tenant struct {
	Name  string        `json:"name"`
	Cloud CloudProvider `json:"cloud"`
	Slug  string        `json:"slug"`
	OrgID string        `json:"org_id"`
	Host  string        `json:"host"`

	KustomizationPath string `json:"kustomizationPath"`

	// ContentPath is where all the tenant resources will be stored
	ContentPath string `json:"contentPath"`
}

type DBCredentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func NewTenant(req TenantRequestBody) (Tenant, error) {
	// TODO: Hardcoded for now
	cloud := Azure

	kPath, err := Template(Config.Git.KustomizationPath, map[string]any{
		"cluster": getClusterName(cloud),
	})
	if err != nil {
		return Tenant{}, err
	}

	slug := req.Data.Slug
	if slug == "" {
		// TODO: If slug is empty, we might have to update the new slug in clerk
		slug = goslug.Make(req.Data.Name)
	}

	return Tenant{
		Name:              req.Data.Name,
		OrgID:             req.Data.OrgID,
		Cloud:             cloud,
		Slug:              slug,
		KustomizationPath: kPath,
		ContentPath:       path.Join(path.Dir(kPath), slug),
		Host:              getHost(cloud, slug),
	}, nil
}

func (tenant Tenant) GenerateDBUsername() string {
	return fmt.Sprintf("%s_%d", strings.ToLower(tenant.Slug), rand.Intn(1000))
}

func (tenant Tenant) GenerateDBPassword() string {
	return RandomString(16)
}

func getClusterName(cloud CloudProvider) string {
	// TODO: Take this from config
	switch cloud {
	case Azure:
		return "azure-internal-prod"
	case AWS:
		return "aws-demo"
	}
	return ""
}

func getHost(cloud CloudProvider, tenantName string) string {
	switch cloud {
	case Azure:
		return fmt.Sprintf("mission-control.%s.internal-prod.flanksource.com", tenantName)
	case AWS:
		return ""
	default:
		return ""
	}
}
