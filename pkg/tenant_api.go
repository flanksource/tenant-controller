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
	AZURE CloudProvider = "azure"
)

type Tenant struct {
	Name  string         `json:"name"`
	Cloud CloudProvider  `json:"cloud"`
	Slug  string         `json:"slug,omitempty"`
	AZURE v1.AZUREConfig `json:"-"`
	DB    DBCredentials  `json:"-"`
}

type DBCredentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (tenant *Tenant) GenerateDBCredentials() *Tenant {
	// Generate a random username and password
	username := fmt.Sprintf("%s_%d", strings.ToLower(tenant.Slug), rand.Intn(1000))
	password := generateRandomPassword()

	// Set the username and password on the tenant object
	tenant.DB.Username = username
	tenant.DB.Password = password

	// Return the updated tenant object
	return tenant
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
