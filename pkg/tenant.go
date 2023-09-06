package pkg

import (
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

const (
	HELM_RELEASE_TEMPLATE = `
apiVersion: helm.toolkit.fluxcd.io/v2beta1
kind: HelmRelease
metadata:
  name: {{.tenant}}
  namespace: {{.tenant}}
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
    domain: {{.tenant}}.internal-prod.flanksource.com
    vcluster:
      syncer:
        extraArgs:
          - --tls-san={{.tenant}}.{{.tenant}}.svc
          - --out-kube-config-server=https://{{.tenant}}.{{.tenant}}.svc
    missionControl:
      authProvider: clerk
      clerkJWKSURL: {{.jwksURL}}
      clerkOrgID: {{.orgID}}
`

	INGRESS_TEMPLATE = `
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  annotations:
    kubernetes.io/tls-acme: "true"
  name: mission-control-{{.tenant}}
  namespace: {{.tenant}}
spec:
  ingressClassName: nginx
  rules:
  - host: mission-control.{{.tenant}}.internal-prod.flanksource.com
    http:
      paths:
      - backend:
          service:
            name: mission-control-x-default-x-{{.tenant}}
            port:
              number: 8080
        path: /
        pathType: ImplementationSpecific
  tls:
  - hosts:
    - mission-control.{{.tenant}}.internal-prod.flanksource.com
    secretName: mission-control-tls-{{.tenant}}
`

	NAMESPACE_TEMPLATE = `
apiVersion: v1
kind: Namespace
metadata:
  name: {{.tenant}}
`

	KUSTOMIZATION_RAW = `
apiVersion: kustomize.config.k8s.io/v1beta1
kind: Kustomization
resources:
  - namespace.yaml
  - helmrelease.yaml
  - secret.yaml
  - ingress.yaml
`
)

func GetTenantResources(tenant Tenant, sealedSecret string) (obj []*unstructured.Unstructured, err error) {
	vars := map[string]any{
		"tenant":  tenant.Slug,
		"jwksURL": Config.Clerk.JWKSURL,
		"orgID":   tenant.OrgID,
	}
	helmReleaseRaw, err := Template(HELM_RELEASE_TEMPLATE, vars)
	if err != nil {
		return nil, err
	}
	namespaceRaw, err := Template(NAMESPACE_TEMPLATE, vars)
	if err != nil {
		return nil, err
	}
	ingressRaw, err := Template(INGRESS_TEMPLATE, vars)
	if err != nil {
		return nil, err
	}

	return GetUnstructuredObjects(namespaceRaw, sealedSecret, KUSTOMIZATION_RAW, helmReleaseRaw, ingressRaw)
}
