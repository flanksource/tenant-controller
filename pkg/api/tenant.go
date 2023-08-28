package api

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/flanksource/commons/logger"
	"github.com/flanksource/tenant-controller/pkg"
	"github.com/flanksource/tenant-controller/pkg/git"
	"github.com/flanksource/tenant-controller/pkg/git/connectors"
	"github.com/gosimple/slug"
	"github.com/labstack/echo/v4"
)

func CreateTenant(c echo.Context) error {
	if c.Request().Body == nil {
		logger.Debugf("missing request body")
		return errorResonse(c, errors.New("missing request body"), http.StatusBadRequest)
	}
	defer c.Request().Body.Close()

	tenant := &pkg.Tenant{}
	reqBody, err := io.ReadAll(c.Request().Body)
	if err != nil {
		return errorResonse(c, err, http.StatusInternalServerError)
	}
	if err := json.Unmarshal(reqBody, tenant); err != nil {
		return errorResonse(c, err, http.StatusBadRequest)
	}

	if tenant.Slug == "" {
		tenant.Slug = slug.Make(tenant.Name)
	}

	tenant.GenerateDBCredentials()

	sc := GetSecretControllerFromCloud(tenant.Cloud)
	sealedSecretRaw, err := sc.GenerateSealedSecret(tenant)
	if err != nil {
		return errorResonse(c, err, http.StatusInternalServerError)
	}

	connector, err := connectors.NewConnector(pkg.Config.GIT)
	if err != nil {
		return errorResonse(c, err, http.StatusInternalServerError)
	}

	objs, err := pkg.GetTenantResources(tenant.Slug, sealedSecretRaw)
	if err != nil {
		return errorResonse(c, err, http.StatusInternalServerError)
	}
	work, title, err := git.CreateTenantResources(connector, tenant, objs)
	if err != nil {
		return errorResonse(c, err, http.StatusInternalServerError)
	}

	hash, err := git.CreateCommit(work, title)
	if err != nil {
		return errorResonse(c, err, http.StatusInternalServerError)
	}
	if err = connector.Push(context.TODO(), fmt.Sprintf("%s:%s", pkg.Config.GIT.Branch, pkg.Config.GIT.Base)); err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	pr, err := connector.OpenPullRequest(context.TODO(), pkg.Config.GIT.Base, pkg.Config.GIT.Branch, pkg.Config.GIT.PullRequest)
	if err != nil {
		return errorResonse(c, err, http.StatusInternalServerError)
	}

	return c.String(http.StatusAccepted, fmt.Sprintf("Committed %s, PR: %d ", hash, pr))
}
