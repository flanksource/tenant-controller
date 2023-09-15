package tenant

import (
	v1 "github.com/flanksource/tenant-controller/api/v1"
	"github.com/flanksource/tenant-controller/pkg/config"
	"github.com/flanksource/tenant-controller/pkg/utils"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

const (
	HELM_RELEASE_TEMPLATE = `
apiVersion: helm.toolkit.fluxcd.io/v2beta1
kind: HelmRelease
metadata:
  name: {{.orgID}}
  namespace: {{.orgID}}
  annotations:
    flanksource.com/tenant-slug: {{.slug}}
spec:
  interval: 5m
  chart:
    spec:
      chart: mission-control-tenant
      sourceRef:
        kind: HelmRepository
        name: flanksource
        namespace: production
      interval: 1m
  install:
    crds: CreateReplace
  upgrade:
    crds: CreateReplace
  values:
    domain: {{.host}}
    vcluster:
      syncer:
        extraArgs:
          - --tls-san={{.orgID}}.{{.orgID}}.svc
          - --out-kube-config-server=https://{{.orgID}}.{{.orgID}}.svc
    missionControl:
      authProvider: clerk
      clerkJWKSURL: {{.jwksURL}}
      clerkOrgID: {{.orgID}}
`

	NAMESPACE_TEMPLATE = `
apiVersion: v1
kind: Namespace
metadata:
  name: {{.orgID}}
`

	KUSTOMIZATION_RAW = `
apiVersion: kustomize.config.k8s.io/v1beta1
kind: Kustomization
resources:
  - namespace.yaml
  - helmrelease.yaml
  - secret.yaml
`
)

func GetTenantResources(tenant v1.Tenant, sealedSecret string) (obj []*unstructured.Unstructured, err error) {
	vars := map[string]any{
		"slug":    tenant.Slug,
		"host":    tenant.Host,
		"jwksURL": config.Config.Clerk.JWKSURL,
		"orgID":   tenant.ID,
	}
	helmReleaseRaw, err := utils.Template(HELM_RELEASE_TEMPLATE, vars)
	if err != nil {
		return nil, err
	}
	namespaceRaw, err := utils.Template(NAMESPACE_TEMPLATE, vars)
	if err != nil {
		return nil, err
	}

	return utils.GetUnstructuredObjects(namespaceRaw, sealedSecret, KUSTOMIZATION_RAW, helmReleaseRaw)
}
