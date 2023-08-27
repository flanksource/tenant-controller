package secrets

import (
	"os"
	"os/exec"

	"github.com/flanksource/tenant-controller/pkg"
)

type AzureSealedSecret struct{}

func (s *AzureSealedSecret) GenerateSealedSecret(tenant *pkg.Tenant) ([]byte, error) {
	fileName, err := createDBSecretFile(tenant.Slug, tenant.DB.Username, tenant.DB.Password)
	if err != nil {
		return nil, err
	}
	pkg.Config.AZURE.SetENVS()
	cmd := exec.Command("sops", "--encrypt", "--encrypted-regex", "stringData", "--azure-kv", os.Getenv("AZURE_VAULT_URL"), fileName)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, err
	}
	return output, nil
}
