#!/bin/bash

# Deploy WentFramework to Kubernetes
# Usage: ./deploy.sh [environment]

set -e

ENVIRONMENT=${1:-"development"}
NAMESPACE="wentframework"

echo "🚀 Deploying WentFramework to Kubernetes..."
echo "📦 Environment: $ENVIRONMENT"
echo "🏷️  Namespace: $NAMESPACE"

# Build Docker image
echo "🔨 Building Docker image..."
docker build -t wentframework:latest .

# Apply Kubernetes manifests
echo "📋 Applying Kubernetes manifests..."

# Create namespace first
kubectl apply -f k8s/namespace.yaml

# Apply all other resources
kubectl apply -f k8s/

# Wait for deployments to be ready
echo "⏳ Waiting for deployments to be ready..."
kubectl wait --for=condition=available --timeout=300s deployment/postgres-deployment -n $NAMESPACE
kubectl wait --for=condition=available --timeout=300s deployment/wentframework-deployment -n $NAMESPACE

# Get service information
echo "✅ Deployment completed successfully!"
echo ""
echo "📊 Service Information:"
kubectl get services -n $NAMESPACE

echo ""
echo "🏃 Pod Status:"
kubectl get pods -n $NAMESPACE

echo ""
echo "🌐 Access Information:"
echo "- Internal Service: wentframework-service.wentframework.svc.cluster.local"
echo "- LoadBalancer: $(kubectl get service wentframework-service -n $NAMESPACE -o jsonpath='{.status.loadBalancer.ingress[0].ip}')"
echo "- Ingress: wentframework.local (add to /etc/hosts if using local setup)"

echo ""
echo "🔧 Useful Commands:"
echo "- View logs: kubectl logs -f deployment/wentframework-deployment -n $NAMESPACE"
echo "- Scale up: kubectl scale deployment/wentframework-deployment --replicas=5 -n $NAMESPACE"
echo "- Port forward: kubectl port-forward service/wentframework-service 8080:80 -n $NAMESPACE"
