package v1

import (
	"os"
)

type CloudProvider string

const (
	AWS   CloudProvider = "aws"
	AZURE CloudProvider = "azure"
)

type Config struct {
	Git   *GitopsAPISpec `json:"git,omitempty" yaml:"git,omitempty"`
	AWS   *AWSConfig     `json:"aws,omitempty" yaml:"aws,omitempty"`
	Azure *AZUREConfig   `json:"azure,omitempty" yaml:"azure,omitempty"`
}

type AWSConfig struct {
	// ARN of the key to use for encryption
	Key string `json:"key" yaml:"key"`
}

type AZUREConfig struct {
	TenantID     string `json:"tenant_id" yaml:"tenant_id"`
	ClientID     string `json:"client_id" yaml:"client_id"`
	ClientSecret string `json:"client_secret" yaml:"client_secret"`
	VaultURI     string `json:"vault_uri" yaml:"vault_url"`
}

func (a *AZUREConfig) SetEnvs() {
	os.Setenv("AZURE_CLIENT_SECRET", a.ClientSecret)
	os.Setenv("AZURE_CLIENT_ID", a.ClientID)
	os.Setenv("AZURE_TENANT_ID", a.TenantID)
	os.Setenv("AZURE_VAULT_URL", a.VaultURI)
}
