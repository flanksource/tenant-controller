package api

import (
	"github.com/flanksource/tenant-controller/pkg"
	"github.com/flanksource/tenant-controller/pkg/secrets"
	"github.com/labstack/echo/v4"
)

func errorResonse(c echo.Context, err error, code int) error {
	e := map[string]string{"error": err.Error()}
	return c.JSON(code, e)
}

func GetSecretControllerFromCloud(cloud pkg.CloudProvider) secrets.Secrets {
	switch cloud {
	case pkg.Azure:
		return &secrets.AzureSealedSecret{}
	}
	return nil
}
