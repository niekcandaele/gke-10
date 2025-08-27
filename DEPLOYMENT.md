# GKE Hackathon Deployment Guide

## Overview
This repository contains Kustomize configurations and Argo CD applications for deploying Bank of Anthos and Online Boutique on GKE.

## Architecture
- **Bank of Anthos**: Banking application demo deployed in `bank-of-anthos` namespace
- **Online Boutique**: E-commerce application demo deployed in `online-boutique` namespace  
- **GitOps**: Argo CD manages deployments from this repository
- **Ingress**: Both applications exposed via NGINX ingress

## Key Modifications
1. **Service Types**: Frontend services changed from LoadBalancer to ClusterIP
2. **Payment Service**: Removed from Online Boutique
3. **Service Accounts**: Fixed StatefulSet service account configuration
4. **JWT Secret**: Added to Bank of Anthos for authentication

## Repository Structure
```
k8s/
├── argocd/
│   └── applications/        # Argo CD Application manifests
├── kustomize/
│   ├── bank-of-anthos/
│   │   ├── base/            # Base configuration with upstream references
│   │   └── overlays/
│   │       └── hackathon/   # Environment-specific patches
│   └── online-boutique/
│       ├── base/            # Base configuration with upstream references
│       └── overlays/
│           └── hackathon/   # Environment-specific patches
```

## Deployment Instructions

### Prerequisites
- GKE cluster with NGINX ingress controller
- Argo CD installed in the cluster
- kubectl configured to access the cluster

### Deploy Applications

1. **Apply Argo CD Applications**:
   ```bash
   kubectl apply -f k8s/argocd/applications/
   ```

2. **Verify Deployment**:
   ```bash
   # Check application status
   kubectl get applications -n argocd
   
   # Check pods
   kubectl get pods -n bank-of-anthos
   kubectl get pods -n online-boutique
   
   # Check ingresses
   kubectl get ingress -A
   ```

3. **Access Applications**:
   - Bank of Anthos: `https://bank.gke10.candaele.dev`
   - Online Boutique: `https://boutique.gke10.candaele.dev`
   
   Note: DNS should be configured with a wildcard A record for `*.gke10.candaele.dev` pointing to the ingress IP. 
   SSL certificates are automatically provisioned by cert-manager using Let's Encrypt.

### Manual Kustomize Build (for testing)

```bash
# Build Bank of Anthos manifests
kubectl kustomize k8s/kustomize/bank-of-anthos/overlays/hackathon/

# Build Online Boutique manifests  
kubectl kustomize k8s/kustomize/online-boutique/overlays/hackathon/
```

## Troubleshooting

### Sync Issues
If Argo CD shows OutOfSync status:
1. Check the application details: `kubectl describe application <app-name> -n argocd`
2. Force sync: `kubectl patch application <app-name> -n argocd --type merge -p '{"operation": {"sync": {}}}'`

### Pod Issues
- Bank of Anthos contacts pod may have Workload Identity issues - this is a GCP IAM configuration issue
- Check pod logs: `kubectl logs <pod-name> -n <namespace>`
- Check events: `kubectl get events -n <namespace> --sort-by='.lastTimestamp'`

## Requirements Validation
- ✅ REQ-001: Bank of Anthos in dedicated namespace
- ✅ REQ-002: Online Boutique in dedicated namespace  
- ✅ REQ-003: Bank frontend Service is ClusterIP
- ✅ REQ-004: Boutique frontend-external removed
- ✅ REQ-005: Ingress resources created
- ✅ REQ-006: Payment service removed from Online Boutique
- ✅ REQ-007: Kustomize overlays working
- ✅ REQ-008: Argo CD managing deployments
- ✅ REQ-009: Frontend deployments running