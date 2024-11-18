package secrets

import (
	"os/exec"

	v1 "github.com/flanksource/tenant-controller/api/v1"
)

type GCPSealedSecret struct{}

func (s GCPSealedSecret) GenerateSealedSecret(params SealedSecretParams) ([]byte, error) {
	fileName, err := createDBSecretFile(params.Namespace, params.Username, params.Password)
	if err != nil {
		return nil, err
	}

	cmd := exec.Command(
		"sops", "--encrypt",
		"--encrypted-regex", "stringData",
		"--gcp-kms", v1.GlobalConfig.GCP.KMS,
		fileName,
	)

	return cmd.CombinedOutput()

}
