package secrets

import (
	"os"

	"github.com/flanksource/tenant-controller/pkg"
	"gopkg.in/yaml.v2"
)

type Secrets interface {
	// GenerateSealedSecret generates a sealed secret from the tenant authenticating using the configured cloud provider
	GenerateSealedSecret(tenant *pkg.Tenant) ([]byte, error)
}

// create a function that creates a kubernetes secret object structure and write it into a file

func createDBSecretFile(namePrefix, username, password string) (string, error) {
	manifest := map[string]interface{}{
		"apiVersion": "v1",
		"kind":       "Secret",
		"metadata": map[string]string{
			"name": namePrefix + "-db-credentials",
		},
		"spec": map[string]interface{}{
			"type": "Opaque",
			"stringData": map[string]string{
				"username": username,
				"password": password,
			},
		},
	}

	// marshal the manifest into YAML
	yamlData, err := yaml.Marshal(manifest)
	if err != nil {
		return "", err
	}

	file, err := os.CreateTemp("", "db-credentials-*.yaml")
	if err != nil {
		return "", err
	}
	defer file.Close()

	_, err = file.Write(yamlData)
	if err != nil {
		return "", err
	}

	// return the file path
	return file.Name(), nil
}