# Kubernetes Deployment for WentFramework

This directory contains Kubernetes manifests and deployment scripts for running WentFramework in a Kubernetes cluster.

## 📁 Files Overview

### Configuration Files
- `namespace.yaml` - Creates the wentframework namespace
- `configmap.yaml` - Application configuration (non-sensitive)
- `secrets.yaml` - Sensitive configuration (passwords, JWT secrets)

### Database
- `postgres-pv.yaml` - Persistent volume and claim for PostgreSQL data
- `postgres-deployment.yaml` - PostgreSQL database deployment
- `postgres-service.yaml` - PostgreSQL internal service

### Application
- `wentframework-deployment.yaml` - Main application deployment with auto-migration
- `wentframework-service.yaml` - LoadBalancer service for external access
- `ingress.yaml` - Ingress configuration for domain-based routing
- `hpa.yaml` - Horizontal Pod Autoscaler for automatic scaling

### Deployment Tools
- `kustomization.yaml` - Kustomize configuration for organized deployments
- `deploy.sh` - Automated deployment script
- `cleanup.sh` - Cleanup script for removing all resources

## 🚀 Quick Deployment

### Prerequisites

1. **Kubernetes cluster** (minikube, kind, GKE, EKS, AKS, etc.)
2. **kubectl** configured to connect to your cluster
3. **Docker** for building the application image

### Deploy Everything

```bash
# Make scripts executable (if not already done)
chmod +x k8s/deploy.sh k8s/cleanup.sh

# Deploy to Kubernetes
./k8s/deploy.sh
```

This will:
- Build the Docker image
- Create all Kubernetes resources
- Wait for deployments to be ready
- Display access information

### Access the Application

After deployment, you can access the application via:

1. **Port Forward** (for testing):
   ```bash
   kubectl port-forward service/wentframework-service 8080:80 -n wentframework
   # Access at: http://localhost:8080
   ```

2. **LoadBalancer** (if supported by your cluster):
   ```bash
   kubectl get service wentframework-service -n wentframework
   # Use the EXTERNAL-IP
   ```

3. **Ingress** (if ingress controller is installed):
   ```bash
   # Add to /etc/hosts: <ingress-ip> wentframework.local
   # Access at: http://wentframework.local
   ```

## 🔧 Manual Deployment

If you prefer to deploy step by step:

```bash
# 1. Build the Docker image
docker build -t wentframework:latest .

# 2. Create namespace
kubectl apply -f k8s/namespace.yaml

# 3. Apply configuration
kubectl apply -f k8s/configmap.yaml
kubectl apply -f k8s/secrets.yaml

# 4. Deploy PostgreSQL
kubectl apply -f k8s/postgres-pv.yaml
kubectl apply -f k8s/postgres-deployment.yaml
kubectl apply -f k8s/postgres-service.yaml

# 5. Deploy application
kubectl apply -f k8s/wentframework-deployment.yaml
kubectl apply -f k8s/wentframework-service.yaml

# 6. Optional: Ingress and autoscaling
kubectl apply -f k8s/ingress.yaml
kubectl apply -f k8s/hpa.yaml
```

## 📊 Monitoring and Management

### Check Status
```bash
# View all resources
kubectl get all -n wentframework

# Check pod logs
kubectl logs -f deployment/wentframework-deployment -n wentframework

# Check database logs
kubectl logs -f deployment/postgres-deployment -n wentframework
```

### Scaling
```bash
# Manual scaling
kubectl scale deployment/wentframework-deployment --replicas=5 -n wentframework

# Auto-scaling is configured via HPA (70% CPU, 80% Memory)
kubectl get hpa -n wentframework
```

### Database Access
```bash
# Connect to PostgreSQL
kubectl exec -it deployment/postgres-deployment -n wentframework -- psql -U went_user -d went_test
```

## 🏗️ Architecture

```
┌─────────────────┐    ┌─────────────────┐
│   Ingress       │    │  LoadBalancer   │
│ wentframework   │    │    Service      │
│     .local      │    │                 │
└─────────┬───────┘    └─────────┬───────┘
          │                      │
          └──────────┬───────────┘
                     │
        ┌────────────▼────────────┐
        │   WentFramework Pods    │
        │    (3 replicas)         │
        │   ┌─────────────────┐   │
        │   │ Init Container  │   │
        │   │ - Wait for DB   │   │
        │   │ - Run Migration │   │
        │   └─────────────────┘   │
        │   ┌─────────────────┐   │
        │   │ Main Container  │   │
        │   │ - Go App :3000  │   │
        │   └─────────────────┘   │
        └─────────┬───────────────┘
                  │
        ┌─────────▼───────────────┐
        │   PostgreSQL Pod        │
        │   - Port: 5432          │
        │   - Persistent Storage  │
        └─────────────────────────┘
```

## 🔒 Security Considerations

### Secrets Management
- Database password is base64 encoded in `secrets.yaml`
- **Production**: Use external secret management (Vault, AWS Secrets Manager, etc.)
- **Change default passwords** before deploying to production

### Network Security
- PostgreSQL is only accessible within the cluster (ClusterIP)
- Application uses LoadBalancer for external access
- Consider using NetworkPolicies for additional isolation

### Resource Limits
- Both PostgreSQL and WentFramework have resource limits configured
- Adjust based on your workload requirements

## 🌍 Environment Configurations

### Development
```bash
# Use with minikube or kind
./k8s/deploy.sh development
```

### Production
1. Update `secrets.yaml` with production credentials
2. Modify resource limits in deployment files
3. Configure ingress with proper TLS certificates
4. Set up monitoring and logging

### Cloud Providers

#### Google GKE
```bash
# Create cluster
gcloud container clusters create wentframework-cluster

# Deploy
./k8s/deploy.sh production
```

#### AWS EKS
```bash
# Create cluster with eksctl
eksctl create cluster --name wentframework-cluster

# Deploy
./k8s/deploy.sh production
```

#### Azure AKS
```bash
# Create cluster
az aks create --resource-group myResourceGroup --name wentframework-cluster

# Deploy
./k8s/deploy.sh production
```

## 🧹 Cleanup

To remove all resources:

```bash
./k8s/cleanup.sh
```

**Note**: This will delete all data. The persistent volume is retained by default.

## 📋 Troubleshooting

### Common Issues

1. **Image Pull Errors**
   ```bash
   # If using local Docker image, load it into your cluster
   # For minikube:
   eval $(minikube docker-env)
   docker build -t wentframework:latest .
   ```

2. **Database Connection Issues**
   ```bash
   # Check if PostgreSQL is running
   kubectl get pods -n wentframework
   kubectl logs deployment/postgres-deployment -n wentframework
   ```

3. **Migration Failures**
   ```bash
   # Check init container logs
   kubectl logs deployment/wentframework-deployment -c migrate -n wentframework
   ```

4. **Service Not Accessible**
   ```bash
   # Check service and endpoints
   kubectl get svc,endpoints -n wentframework
   ```

### Debug Commands
```bash
# Get detailed pod information
kubectl describe pod <pod-name> -n wentframework

# Execute into running container
kubectl exec -it <pod-name> -n wentframework -- /bin/sh

# View events
kubectl get events -n wentframework --sort-by=.metadata.creationTimestamp
```

## 🔄 Updates and Rollbacks

### Rolling Updates
```bash
# Update image tag
kubectl set image deployment/wentframework-deployment wentframework=wentframework:v2.0.0 -n wentframework

# Monitor rollout
kubectl rollout status deployment/wentframework-deployment -n wentframework
```

### Rollbacks
```bash
# View rollout history
kubectl rollout history deployment/wentframework-deployment -n wentframework

# Rollback to previous version
kubectl rollout undo deployment/wentframework-deployment -n wentframework
```

---

**Happy Kubernetes deployment! 🚀**
