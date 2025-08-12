#!/bin/bash

# Cleanup WentFramework from Kubernetes
# Usage: ./cleanup.sh

set -e

NAMESPACE="wentframework"

echo "🧹 Cleaning up WentFramework from Kubernetes..."
echo "🏷️  Namespace: $NAMESPACE"

# Delete all resources
echo "🗑️  Deleting Kubernetes resources..."
kubectl delete -f k8s/ --ignore-not-found=true

# Delete persistent volumes (optional - uncomment if you want to delete data)
# echo "💾 Deleting persistent volumes..."
# kubectl delete pv postgres-pv --ignore-not-found=true

echo "✅ Cleanup completed successfully!"
echo "🔍 Remaining resources (should be empty):"
kubectl get all -n $NAMESPACE 2>/dev/null || echo "Namespace $NAMESPACE not found (expected)"
