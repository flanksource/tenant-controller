package tenant

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/flanksource/commons/logger"
	v1 "github.com/flanksource/tenant-controller/api/v1"
	"github.com/flanksource/tenant-controller/pkg/git"
	"github.com/flanksource/tenant-controller/pkg/git/connectors"
	"github.com/flanksource/tenant-controller/pkg/secrets"
	"github.com/flanksource/tenant-controller/pkg/utils"
	"github.com/labstack/echo/v4"
	"github.com/samber/lo"
	k8sv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
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

	if reqBody.Type != "organization.created" {
		return c.String(http.StatusAlreadyReported, "Only organization.created events are accepted")
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
		return errorResonse(c, fmt.Errorf("error creating tenant resources: %w", err), http.StatusInternalServerError)
	}

	pr, hash, err := git.OpenPRWithTenantResources(tenant, objs)
	if err != nil {
		return errorResonse(c, fmt.Errorf("erorr commiting to git: %w", err), http.StatusInternalServerError)
	}

	return c.String(http.StatusAccepted, fmt.Sprintf("Committed %s, PR: %d ", hash, pr))
}

func Reconcile(k8s kubernetes.Interface) error {
	kPath, err := utils.Template(v1.GlobalConfig.Git.KustomizationPath, map[string]any{
		"cluster": v1.GlobalConfig.GetClusterName(),
	})
	if err != nil {
		return err
	}

	connector, err := connectors.NewConnector(v1.GlobalConfig.Git)
	if err != nil {
		return err
	}

	fs, _, err := connector.Clone(context.Background(), v1.GlobalConfig.Git.PullRequest.Base, v1.GlobalConfig.Git.PullRequest.Base)
	if err != nil {
		return err
	}

	kust, err := git.GetKustomizaton(fs, kPath)
	if err != nil {
		return err
	}
	orgsInKustomize := lo.Filter(kust.Resources, func(r string, _ int) bool { return strings.HasPrefix(r, "org-") })

	nsList, err := k8s.CoreV1().Namespaces().List(context.Background(), metav1.ListOptions{})
	if err != nil {
		return err
	}

	orgsInCluster := lo.Filter(lo.Map(nsList.Items, func(ns k8sv1.Namespace, _ int) string {
		return ns.Name
	}), func(ns string, _ int) bool {
		return strings.HasPrefix(ns, "org-")
	})

	orgsToRemove, _ := lo.Difference(orgsInCluster, orgsInKustomize)

	for _, org := range orgsToRemove {
		if err := k8s.CoreV1().Namespaces().Delete(context.Background(), org, metav1.DeleteOptions{}); err != nil {
			logger.Errorf("error deleting namespace[%s]: %v", org, err)
		}
		logger.Infof("Deleted %s", org)
	}

	return nil
}
