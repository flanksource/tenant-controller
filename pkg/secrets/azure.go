package secrets

import (
	"fmt"
	"os/exec"

	"github.com/flanksource/tenant-controller/pkg/config"
)

type AzureSealedSecret struct{}

func (s *AzureSealedSecret) GenerateSealedSecret(params SealedSecretParams) ([]byte, error) {
	fileName, err := createDBSecretFile(params.Namespace, params.Username, params.Password)
	if err != nil {
		return nil, err
	}

	cmd := exec.Command(
		"sops", "--encrypt",
		"--encrypted-regex", "stringData",
		"--azure-kv", config.Config.Azure.VaultURI,
		fileName,
	)

	cmd.Env = append(cmd.Env,
		fmt.Sprintf("AZURE_CLIENT_ID=%s", config.Config.Azure.ClientID),
		fmt.Sprintf("AZURE_TENANT_ID=%s", config.Config.Azure.TenantID),
		fmt.Sprintf("AZURE_CLIENT_SECRET=%s", config.Config.Azure.ClientSecret),
	)

	return cmd.CombinedOutput()
}
