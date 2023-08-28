package secrets

import (
	"os"
	"os/exec"

	"github.com/flanksource/tenant-controller/pkg"
)

type AzureSealedSecret struct{}

func (s *AzureSealedSecret) GenerateSealedSecret(params SealedSecretParams) ([]byte, error) {
	fileName, err := createDBSecretFile(params.Slug, params.Username, params.Password)
	if err != nil {
		return nil, err
	}
	pkg.Config.Azure.SetEnvs()
	cmd := exec.Command("sops", "--encrypt", "--encrypted-regex", "stringData", "--azure-kv", os.Getenv("AZURE_VAULT_URL"), fileName)
	return cmd.CombinedOutput()
}
