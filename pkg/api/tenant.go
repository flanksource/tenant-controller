package api

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/flanksource/tenant-controller/pkg"
	"github.com/flanksource/tenant-controller/pkg/git"
	"github.com/gosimple/slug"
	"github.com/labstack/echo/v4"
)

func CreateTenant(c echo.Context) error {
	if c.Request().Body == nil {
		return errorResonse(c, errors.New("missing request body"), http.StatusBadRequest)
	}
	defer c.Request().Body.Close()

	var tenant *pkg.Tenant
	if err := c.Bind(tenant); err != nil {
		return errorResonse(c, err, http.StatusBadRequest)
	}

	if tenant.Slug == "" {
		tenant.Slug = slug.Make(tenant.Name)
	}

	// TODO: Remove side-effects
	// This sets certain parameters that are used
	// by sc.GenrateSealedSecret
	// Use `SealedSecretParams` struct to pass into that func
	tenant.GenerateDBCredentials()

	// TODO: Webhook does not tell which cloud provider
	sc := GetSecretControllerFromCloud(tenant.Cloud)
	sealedSecretRaw, err := sc.GenerateSealedSecret(tenant)
	if err != nil {
		return errorResonse(c, err, http.StatusInternalServerError)
	}

	objs, err := pkg.GetTenantResources(tenant.Slug, sealedSecretRaw)
	if err != nil {
		return errorResonse(c, err, http.StatusInternalServerError)
	}

	pr, hash, err := git.OpenPRWithTenantResources(tenant, objs)
	if err != nil {
		return errorResonse(c, err, http.StatusInternalServerError)
	}

	return c.String(http.StatusAccepted, fmt.Sprintf("Committed %s, PR: %d ", hash, pr))
}
