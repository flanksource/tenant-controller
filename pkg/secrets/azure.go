package secrets

import (
	"os"
	"os/exec"

	"github.com/flanksource/tenant-controller/pkg/config"
)

type AzureSealedSecret struct{}

func (s *AzureSealedSecret) GenerateSealedSecret(params SealedSecretParams) ([]byte, error) {
	fileName, err := createDBSecretFile(params.Namespace, params.Username, params.Password)
	if err != nil {
		return nil, err
	}
	config.Config.Azure.SetEnvs()

	return exec.Command(
		"sops", "--encrypt",
		"--encrypted-regex", "stringData",
		"--azure-kv", os.Getenv("AZURE_VAULT_URL"),
		fileName,
	).CombinedOutput()
}
