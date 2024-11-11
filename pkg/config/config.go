package config

import (
	"fmt"
	"os"
	"strings"

	yamlutil "k8s.io/apimachinery/pkg/util/yaml"

	v1 "github.com/flanksource/tenant-controller/api/v1"
)

func SetConfig(configFile string) error {
	config, err := ParseConfig(configFile)
	if err != nil {
		return err
	}
	if err := validateConfig(config); err != nil {
		return err
	}
	v1.GlobalConfig = config
	return nil
}

func ParseConfig(configFile string) (config *v1.Config, err error) {
	rawConfig, err := readFile(configFile)
	if err != nil {
		return nil, err
	}
	decoder := yamlutil.NewYAMLOrJSONDecoder(strings.NewReader(rawConfig), 1024)

	if err := decoder.Decode(&config); err != nil {
		return nil, err
	}
	return
}

// Certain fields have to be set so check
// those at start time and panic if validation fails
func validateConfig(config *v1.Config) error {
	if config.DefaultCloud == "" {
		return fmt.Errorf("default_cloud cannot be empty")
	}

	if config.DefaultCloud == v1.Azure {
		if config.Azure.TenantHostFormat == "" {
			return fmt.Errorf("azure.tenant_host_fmt cannot be empty")
		}
		if config.Azure.TenantCluster == "" {
			return fmt.Errorf("azure.tenant_cluster cannot be empty")
		}
	}
	return nil
}

func readFile(path string) (str string, err error) {
	var data []byte
	if data, err = os.ReadFile(path); err != nil {
		return "", err
	}
	return string(data), nil
}
