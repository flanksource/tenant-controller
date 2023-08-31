package pkg

import (
	"fmt"
	"math/rand"
	"path"
	"strings"
)

type CloudProvider string

const (
	AWS   CloudProvider = "aws"
	Azure CloudProvider = "azure"
)

type TenantRequestBody struct {
	Name  string        `json:"name"`
	Cloud CloudProvider `json:"cloud"`
	Slug  string        `json:"slug,omitempty"`
}

type Tenant struct {
	Name  string        `json:"name"`
	Cloud CloudProvider `json:"cloud"`
	Slug  string        `json:"slug,omitempty"`

	// Not sure why this was added
	// But commenting out since it is not in use
	//Azure             v1.AzureConfig `json:"-"`

	KustomizationPath string `json:"kustomizationPath"`

	// ContentPath is where all the tenant resources will be stored
	ContentPath string `json:"contentPath"`
}

type DBCredentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func NewTenant(t TenantRequestBody) (Tenant, error) {
	kPath, err := Template(Config.Git.KustomizationPath, map[string]any{
		"cluster": getClusterName(t.Cloud),
	})
	if err != nil {
		return Tenant{}, err
	}

	contentPath := path.Join(path.Dir(kPath), t.Slug)

	return Tenant{
		Name:              t.Name,
		Cloud:             t.Cloud,
		Slug:              t.Slug,
		KustomizationPath: kPath,
		ContentPath:       contentPath,
	}, nil
}

func (tenant Tenant) GenerateDBUsername() string {
	return fmt.Sprintf("%s_%d", strings.ToLower(tenant.Slug), rand.Intn(1000))
}

func (tenant Tenant) GenerateDBPassword() string {
	return generateRandomPassword()
}

func generateRandomPassword() string {
	// Generate a random password of length 16
	const passwordLength = 16
	const passwordChars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!@#$%^&*()-_=+[]{}|;:,.<>/?"
	password := make([]byte, passwordLength)
	for i := range password {
		password[i] = passwordChars[rand.Intn(len(passwordChars))]
	}
	return string(password)
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
