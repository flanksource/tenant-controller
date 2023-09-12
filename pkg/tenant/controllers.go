package tenant

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/flanksource/commons/logger"
	v1 "github.com/flanksource/tenant-controller/api/v1"
	"github.com/flanksource/tenant-controller/pkg/git"
	"github.com/flanksource/tenant-controller/pkg/secrets"
	"github.com/labstack/echo/v4"
)

func CreateTenant(c echo.Context) error {
	if c.Request().Body == nil {
		return errorResonse(c, errors.New("missing request body"), http.StatusBadRequest)
	}
	defer c.Request().Body.Close()

	var reqBody v1.TenantRequestBody
	if err := c.Bind(&reqBody); err != nil {
		logger.Infof("Broken %v", err)
		return errorResonse(c, err, http.StatusBadRequest)
	}

	tenant, err := NewTenant(reqBody)
	if err != nil {
		return errorResonse(c, err, http.StatusBadRequest)
	}

	if err := updateHostOnClerk(tenant.OrgID, tenant.Host); err != nil {
		return errorResonse(c, err, http.StatusInternalServerError)
	}

	// TODO: Webhook does not tell which cloud provider
	sc := GetSecretControllerFromCloud(tenant.Cloud)
	sealedSecretRaw, err := sc.GenerateSealedSecret(secrets.SealedSecretParams{
		Slug:     tenant.Slug,
		Username: tenant.DBUsername,
		Password: tenant.DBPassword,
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
