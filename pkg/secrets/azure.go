package secrets

import (
	"fmt"
	"os/exec"

	"github.com/flanksource/tenant-controller/api/v1"
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
		"--azure-kv", v1.GlobalConfig.Azure.VaultURI,
		fileName,
	)

	// If azure.CLIENT_ID is set in config, then add all azure config variables
	// to cmd env, the fallback and desired method is to use workload identity
	if v1.GlobalConfig.Azure.ClientID != "" {
		cmd.Env = append(cmd.Env,
			fmt.Sprintf("AZURE_CLIENT_ID=%s", v1.GlobalConfig.Azure.ClientID),
			fmt.Sprintf("AZURE_TENANT_ID=%s", v1.GlobalConfig.Azure.TenantID),
			fmt.Sprintf("AZURE_CLIENT_SECRET=%s", v1.GlobalConfig.Azure.ClientSecret),
		)
	}

	return cmd.CombinedOutput()
}
