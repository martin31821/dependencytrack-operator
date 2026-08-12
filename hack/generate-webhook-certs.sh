#!/usr/bin/env bash
# Generate self-signed webhook certificates for the DependencyTrack operator.
# The certs are used by both kustomize (config/certs/) and the Helm chart
# (deploy/charts/dependencytrack-operator/certs/). They are generated at
# build time and must NOT be committed to git.
set -euo pipefail

readonly REPO_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
readonly KUSTOMIZE_CERTS_DIR="${REPO_ROOT}/config/certs"
readonly HELM_CERTS_DIR="${REPO_ROOT}/deploy/charts/dependencytrack-operator/certs"
readonly WEBHOOK_MANIFESTS="${REPO_ROOT}/config/webhook/manifests.yaml"
readonly CN="deptrack-operator-webhook"
readonly KUBEWIZED_NAME="deptrack-operator-webhook-service"
readonly KUBEWIZED_NAMESPACE="deptrack-operator-system"

mkdir -p "${KUSTOMIZE_CERTS_DIR}" "${HELM_CERTS_DIR}"

# Build the SAN list with DNS. prefix
SAN_DNS_LIST=""
for san in \
  "deptrack-operator-controller-manager.deptrack-operator-system.svc" \
  "deptrack-operator-controller-manager.deptrack-operator-system.svc.cluster.local" \
  "deptrack-operator-controller-manager-webhook.deptrack-operator-system.svc" \
  "deptrack-operator-controller-manager-webhook.deptrack-operator-system.svc.cluster.local"; do
  if [ -z "$SAN_DNS_LIST" ]; then
    SAN_DNS_LIST="DNS:${san}"
  else
    SAN_DNS_LIST="${SAN_DNS_LIST},DNS:${san}"
  fi
done

# Generate the self-signed cert (10-year validity)
openssl req -x509 -newkey rsa:2048 -nodes \
  -keyout "${KUSTOMIZE_CERTS_DIR}/tls.key" \
  -out  "${KUSTOMIZE_CERTS_DIR}/tls.crt" \
  -days 3650 \
  -subj "/CN=${CN}" \
  -addext "subjectAltName = ${SAN_DNS_LIST}" \
  2>/dev/null

# Copy to Helm chart directory (with Helm's file names)
cp "${KUSTOMIZE_CERTS_DIR}/tls.crt" "${HELM_CERTS_DIR}/webhook.pem"
cp "${KUSTOMIZE_CERTS_DIR}/tls.key" "${HELM_CERTS_DIR}/webhook.key"

# Extract the base64-encoded cert for the caBundle
CA_BUNDLE=$(base64 -w0 < "${KUSTOMIZE_CERTS_DIR}/tls.crt")

# Generate the webhook manifests with the caBundle injected
# The service name and namespace are the kustomize-transformed values
# (namePrefix + namespace applied by kustomize at build time)
cat > "${WEBHOOK_MANIFESTS}" << EOF
---
apiVersion: admissionregistration.k8s.io/v1
kind: ValidatingWebhookConfiguration
metadata:
  name: validating-webhook-configuration
webhooks:
- admissionReviewVersions:
  - v1
  clientConfig:
    service:
      name: ${KUBEWIZED_NAME}
      namespace: ${KUBEWIZED_NAMESPACE}
      path: /validate-dependencytrack-mko-dev-v1alpha1-team
    caBundle: ${CA_BUNDLE}
  failurePolicy: Fail
  name: team.validating.dependencytrack.mko.dev
  rules:
  - apiGroups:
    - dependencytrack.mko.dev
    apiVersions:
    - v1alpha1
    operations:
    - CREATE
    - UPDATE
    resources:
    - teams
  sideEffects: None
EOF

echo "Generated webhook certificates and manifests:"
echo "  Kustomize: ${KUSTOMIZE_CERTS_DIR}/tls.crt, ${KUSTOMIZE_CERTS_DIR}/tls.key"
echo "  Helm:      ${HELM_CERTS_DIR}/webhook.pem, ${HELM_CERTS_DIR}/webhook.key"
echo "  Manifests: ${WEBHOOK_MANIFESTS} (caBundle injected)"