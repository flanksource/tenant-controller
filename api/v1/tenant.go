package v1

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
	ID    string        `json:"id"`
	OrgID string        `json:"org_id"`
	Host  string        `json:"host"`

	KustomizationPath string `json:"kustomizationPath"`

	// ContentPath is where all the tenant resources will be stored
	ContentPath string `json:"contentPath"`

	DBUsername string
	DBPassword string
}

type DBCredentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
