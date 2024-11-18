package secrets

type GCPSealedSecret struct{}

func (s GCPSealedSecret) GenerateSealedSecret(params SealedSecretParams) ([]byte, error) {
	// We setup SQLUser and database via CRDs
	return []byte(""), nil
}
