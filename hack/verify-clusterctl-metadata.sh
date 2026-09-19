#!/usr/bin/env bash
set -euo pipefail

metadata="${1:-metadata.yaml}"

[[ -f "$metadata" ]] || { echo "missing metadata file: $metadata" >&2; exit 1; }
grep -qx 'apiVersion: clusterctl.cluster.x-k8s.io/v1alpha3' "$metadata" || {
  echo 'invalid clusterctl metadata apiVersion' >&2; exit 1;
}
grep -qx 'kind: Metadata' "$metadata" || {
  echo 'clusterctl >=1.11 requires kind: Metadata' >&2; exit 1;
}
grep -qE '^[[:space:]]+releaseSeries:' "$metadata" || {
  echo 'metadata releaseSeries is missing' >&2; exit 1;
}
grep -qE '^[[:space:]]+contract: v1beta1$' "$metadata" || {
  echo 'CAPONE 0.1.x must declare the v1beta1 provider contract' >&2; exit 1;
}

echo 'clusterctl metadata contract: OK'
