package pkg

import (
	"fmt"
	"math/rand"
	"strings"

	v1 "github.com/flanksource/tenant-controller/api/v1"
)

type CloudProvider string

const (
	AWS   CloudProvider = "aws"
	Azure CloudProvider = "azure"
)

type Tenant struct {
	Name  string         `json:"name"`
	Cloud CloudProvider  `json:"cloud"`
	Slug  string         `json:"slug,omitempty"`
	Azure v1.AzureConfig `json:"-"`
}

type DBCredentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
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
