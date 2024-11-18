package v1

import "fmt"

type CloudProvider string

const (
	AWS   CloudProvider = "aws"
	Azure CloudProvider = "azure"
	GCP   CloudProvider = "gcp"
)

var GlobalConfig *Config

func (cloud CloudProvider) GetClusterName() string {
	// TODO: Take this from config
	switch cloud {
	case Azure:
		return GlobalConfig.Azure.TenantCluster
	case AWS:
		return GlobalConfig.AWS.TenantCluster
	case GCP:
		return GlobalConfig.GCP.TenantCluster
	}
	return ""
}

func (cloud CloudProvider) GetHost(tenantID string) string {
	switch cloud {
	case Azure:
		return fmt.Sprintf(GlobalConfig.Azure.TenantHostFormat, tenantID)
	case AWS:
		return fmt.Sprintf(GlobalConfig.AWS.TenantHostFormat, tenantID)
	case GCP:
		return fmt.Sprintf(GlobalConfig.GCP.TenantHostFormat, tenantID)
	default:
		return ""
	}
}

type Config struct {
	DefaultCloud CloudProvider  `json:"default_cloud" yaml:"default_cloud"`
	Git          *GitopsAPISpec `json:"git" yaml:"git"`
	AWS          *AWSConfig     `json:"aws" yaml:"aws"`
	Azure        *AzureConfig   `json:"azure" yaml:"azure"`
	GCP          *GCPConfig     `json:"gcp" yaml:"azure"`
	Clerk        ClerkConfig    `json:"clerk" yaml:"clerk"`
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

type GCPConfig struct {
	KMS              string `json:"kms" yaml:"kms"`
	TenantCluster    string `json:"tenant_cluster" yaml:"tenant_cluster"`
	TenantHostFormat string `json:"tenant_host_fmt" yaml:"tenant_host_fmt"`
}

type ClerkConfig struct {
	SecretKey     string `json:"secretKey" yaml:"secretKey"`
	JWKSURL       string `json:"jwks_url" yaml:"jwks_url"`
	WebhookSecret string `json:"webhook_secret" yaml:"webhook_secret"`
}

func (c Config) GetClusterName() string {
	return c.DefaultCloud.GetClusterName()
}

func (c Config) GetHost(tenantID string) string {
	return c.DefaultCloud.GetHost(tenantID)
}
