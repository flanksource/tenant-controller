package tenant

import (
	v1 "github.com/flanksource/tenant-controller/api/v1"
	"github.com/flanksource/tenant-controller/pkg/secrets"
	"github.com/labstack/echo/v4"
)

func errorResonse(c echo.Context, err error, code int) error {
	e := map[string]string{"error": err.Error()}
	return c.JSON(code, e)
}

func GetSecretControllerFromCloud(cloud v1.CloudProvider) secrets.Secrets {
	switch cloud {
	case v1.Azure:
		return &secrets.AzureSealedSecret{}
	case v1.GCP:
		return &secrets.GCPSealedSecret{}
	}
	return nil
}
