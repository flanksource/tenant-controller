package tenant

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	v1 "github.com/flanksource/tenant-controller/api/v1"
	"github.com/flanksource/tenant-controller/pkg/git"
	"github.com/flanksource/tenant-controller/pkg/secrets"
	"github.com/labstack/echo/v4"
)

var ClerkTenantWebhook *Webhook

func CreateTenant(c echo.Context) error {
	if c.Request().Body == nil {
		return errorResonse(c, errors.New("missing request body"), http.StatusBadRequest)
	}
	defer c.Request().Body.Close()

	body, err := io.ReadAll(c.Request().Body)
	if err != nil {
		errorResonse(c, err, http.StatusBadRequest)
	}

	var reqBody v1.TenantRequestBody
	if err := json.Unmarshal(body, &reqBody); err != nil {
		return errorResonse(c, err, http.StatusBadRequest)
	}

	// Ignoring timestamp since the tolerance is 5mins
	// How to replay older message for whom tenant creation failed ?
	if err := ClerkTenantWebhook.VerifyIgnoringTimestamp(body, c.Request().Header); err != nil {
		return errorResonse(c, fmt.Errorf("webhook verification failed: %w", err), http.StatusBadRequest)
	}

	tenant, err := NewTenant(reqBody)
	if err != nil {
		return errorResonse(c, err, http.StatusBadRequest)
	}

	if err := updateParamsOnClerk(tenant); err != nil {
		return errorResonse(c, err, http.StatusInternalServerError)
	}

	// TODO: Webhook does not tell which cloud provider
	sc := GetSecretControllerFromCloud(tenant.Cloud)
	sealedSecretRaw, err := sc.GenerateSealedSecret(secrets.SealedSecretParams{
		Namespace: tenant.ID,
		Username:  tenant.DBUsername,
		Password:  tenant.DBPassword,
	})
	if err != nil {
		return errorResonse(c, fmt.Errorf("Error generating sealed secret: %s %v", string(sealedSecretRaw), err), http.StatusInternalServerError)
	}

	objs, err := GetTenantResources(tenant, string(sealedSecretRaw))
	if err != nil {
		return errorResonse(c, err, http.StatusInternalServerError)
	}

	pr, hash, err := git.OpenPRWithTenantResources(tenant, objs)
	if err != nil {
		return errorResonse(c, err, http.StatusInternalServerError)
	}

	return c.String(http.StatusAccepted, fmt.Sprintf("Committed %s, PR: %d ", hash, pr))
}
