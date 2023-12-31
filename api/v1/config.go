package v1

type CloudProvider string

const (
	AWS   CloudProvider = "aws"
	AZURE CloudProvider = "azure"
)

type Config struct {
	Git   *GitopsAPISpec `json:"git" yaml:"git"`
	AWS   *AWSConfig     `json:"aws" yaml:"aws"`
	Azure *AzureConfig   `json:"azure" yaml:"azure"`
	Clerk ClerkConfig    `json:"clerk" yaml:"clerk"`
}

type AWSConfig struct {
	// ARN of the key to use for encryption
	Key              string `json:"key" yaml:"key"`
	TenantCluster    string `json:"tenant_cluster" yaml:"tenant_cluster"`
	TenantHostFormat string `json:"tenant_host_fmt" yaml:"tenant_host_fmt"`
}

type AzureConfig struct {
	TenantID         string `json:"tenant_id" yaml:"tenant_id"`
	ClientID         string `json:"client_id" yaml:"client_id"`
	ClientSecret     string `json:"client_secret" yaml:"client_secret"`
	VaultURI         string `json:"vault_uri" yaml:"vault_url"`
	TenantCluster    string `json:"tenant_cluster" yaml:"tenant_cluster"`
	TenantHostFormat string `json:"tenant_host_fmt" yaml:"tenant_host_fmt"`
}

type ClerkConfig struct {
	SecretKey     string `json:"secretKey" yaml:"secretKey"`
	JWKSURL       string `json:"jwks_url" yaml:"jwks_url"`
	WebhookSecret string `json:"webhook_secret" yaml:"webhook_secret"`
}
