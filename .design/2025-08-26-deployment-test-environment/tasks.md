# Implementation Tasks: GKE Hackathon Test Environment Deployment

## Overview
We're building a Kustomize-based deployment configuration for Bank of Anthos and Online Boutique applications on GKE, modifying their service exposure from LoadBalancer to Ingress, and removing the payment service from Boutique. The deployment will be managed via Argo CD for GitOps.

This implementation is divided into 5 phases to incrementally build and test each component.

## Phase 1: Kustomize Base Structure Setup
**Goal**: Create the foundational Kustomize directory structure and base configurations
**Demo**: "At standup, I can show: the Kustomize base structure that references upstream manifests and successfully builds with `kubectl kustomize`"

### Tasks
- [x] Task 1.1: Create Kustomize directory structure
  - **Output**: Complete directory tree for both applications
  - **Files**: 
    - `kustomize/bank-of-anthos/base/`
    - `kustomize/bank-of-anthos/overlays/hackathon/`
    - `kustomize/online-boutique/base/`
    - `kustomize/online-boutique/overlays/hackathon/`
  - **Verify**: `tree kustomize/` shows correct structure

- [x] Task 1.2: Create Bank of Anthos base kustomization
  - **Depends on**: 1.1
  - **Output**: Base kustomization referencing upstream manifests
  - **Files**: `kustomize/bank-of-anthos/base/kustomization.yaml`
  - **Verify**: `kubectl kustomize kustomize/bank-of-anthos/base/` outputs valid YAML

- [x] Task 1.3: Create Online Boutique base kustomization
  - **Depends on**: 1.1
  - **Output**: Base kustomization referencing upstream manifests
  - **Files**: `kustomize/online-boutique/base/kustomization.yaml`
  - **Verify**: `kubectl kustomize kustomize/online-boutique/base/` outputs valid YAML

### Phase 1 Checkpoint
- [x] Run validation: `kubectl kustomize kustomize/bank-of-anthos/base/ | kubectl apply --dry-run=client -f -`
- [x] Run validation: `kubectl kustomize kustomize/online-boutique/base/ | kubectl apply --dry-run=client -f -`
- [x] Manual verification: Both base configurations reference all necessary upstream manifests
- [x] **Demo ready**: Show kustomize build output for both base configurations

## Phase 2: Service Patches and Removal
**Goal**: Create patches to modify service types and remove payment service
**Demo**: "At standup, I can show: Kustomize patches that change LoadBalancer to ClusterIP and remove the payment service"

### Tasks
- [x] Task 2.1: Create Bank of Anthos frontend service patch
  - **Output**: Patch to change frontend Service type to ClusterIP
  - **Files**: `kustomize/bank-of-anthos/overlays/hackathon/patches/frontend-service-patch.yaml`
  - **Verify**: Patch changes `type: LoadBalancer` to `type: ClusterIP`

- [x] Task 2.2: Create Online Boutique frontend service patches
  - **Output**: Patches for frontend-external service modification
  - **Files**: `kustomize/online-boutique/overlays/hackathon/patches/frontend-service-patch.yaml`
  - **Verify**: Patch removes or modifies frontend-external service

- [x] Task 2.3: Create payment service removal patch
  - **Output**: Patch to remove payment service completely
  - **Files**: `kustomize/online-boutique/overlays/hackathon/patches/remove-payment.yaml`
  - **Verify**: Patch removes both Deployment and Service for paymentservice

- [x] Task 2.4: Create overlay kustomization files
  - **Depends on**: 2.1, 2.2, 2.3
  - **Output**: Overlay configurations that apply patches
  - **Files**: 
    - `kustomize/bank-of-anthos/overlays/hackathon/kustomization.yaml`
    - `kustomize/online-boutique/overlays/hackathon/kustomization.yaml`
  - **Verify**: `kubectl kustomize kustomize/*/overlays/hackathon/` shows patched resources

### Phase 2 Checkpoint
- [x] Validate Bank patches: `kubectl kustomize kustomize/bank-of-anthos/overlays/hackathon/ | grep -A5 "name: frontend"`
- [x] Validate Boutique patches: `kubectl kustomize kustomize/online-boutique/overlays/hackathon/ | grep -c paymentservice` returns 0
- [x] Manual verification: Frontend services show ClusterIP, payment service absent
- [x] **Demo ready**: Show diff between base and overlay outputs

## Phase 3: Ingress Configuration
**Goal**: Add Ingress resources to expose frontend services externally
**Demo**: "At standup, I can show: Ingress configurations routing to frontend services with proper hostnames"

### Tasks
- [x] Task 3.1: Create Bank of Anthos Ingress resource
  - **Output**: Ingress configuration for Bank of Anthos
  - **Files**: `kustomize/bank-of-anthos/overlays/hackathon/resources/ingress.yaml`
  - **Verify**: Ingress points to frontend service on port 80

- [x] Task 3.2: Create Online Boutique Ingress resource
  - **Output**: Ingress configuration for Online Boutique
  - **Files**: `kustomize/online-boutique/overlays/hackathon/resources/ingress.yaml`
  - **Verify**: Ingress points to frontend service on port 80

- [x] Task 3.3: Update overlay kustomizations to include Ingress
  - **Depends on**: 3.1, 3.2
  - **Output**: Updated overlays including Ingress resources
  - **Files**: 
    - `kustomize/bank-of-anthos/overlays/hackathon/kustomization.yaml`
    - `kustomize/online-boutique/overlays/hackathon/kustomization.yaml`
  - **Verify**: `kubectl kustomize` output includes Ingress resources

### Phase 3 Checkpoint
- [x] Validate Ingress: `kubectl kustomize kustomize/bank-of-anthos/overlays/hackathon/ | grep -A10 "kind: Ingress"`
- [x] Validate Ingress: `kubectl kustomize kustomize/online-boutique/overlays/hackathon/ | grep -A10 "kind: Ingress"`
- [x] Manual verification: Both Ingress resources have correct hostnames and backends
- [x] **Demo ready**: Show complete Kustomize output with services and ingresses

## Phase 4: Argo CD Applications
**Goal**: Create Argo CD Application manifests for GitOps deployment
**Demo**: "At standup, I can show: Argo CD Application manifests ready for deployment"

### Tasks
- [x] Task 4.1: Create Argo CD applications directory
  - **Output**: Directory for Argo CD manifests
  - **Files**: `argocd/applications/`
  - **Verify**: Directory exists

- [x] Task 4.2: Create Bank of Anthos Argo CD Application
  - **Depends on**: 4.1
  - **Output**: Argo CD Application manifest for Bank of Anthos
  - **Files**: `argocd/applications/bank-of-anthos.yaml`
  - **Verify**: Valid Application resource with correct source and destination

- [x] Task 4.3: Create Online Boutique Argo CD Application
  - **Depends on**: 4.1
  - **Output**: Argo CD Application manifest for Online Boutique
  - **Files**: `argocd/applications/online-boutique.yaml`
  - **Verify**: Valid Application resource with correct source and destination

- [x] Task 4.4: Add namespace configurations
  - **Output**: Namespace definitions in overlays
  - **Files**: Update overlay kustomization files to include namespace
  - **Verify**: Kustomize output includes namespace specifications

### Phase 4 Checkpoint
- [x] Validate Applications: `kubectl apply --dry-run=client -f argocd/applications/`
- [x] Manual verification: Applications point to correct Git paths and namespaces
- [x] **Demo ready**: Show Argo CD Application configurations with sync policies

## Phase 5: Deployment and Validation
**Goal**: Deploy to GKE cluster and validate all functionality
**Demo**: "At standup, I can show: Both applications running on GKE with working Ingress access"

### Tasks
- [ ] Task 5.1: Create namespaces in cluster
  - **Output**: Namespaces created in GKE
  - **Commands**: 
    - `kubectl create namespace bank-of-anthos`
    - `kubectl create namespace online-boutique`
  - **Verify**: `kubectl get namespaces` shows both namespaces

- [ ] Task 5.2: Deploy Bank of Anthos via kubectl (pre-Argo test)
  - **Depends on**: 5.1
  - **Output**: Bank of Anthos running in cluster
  - **Commands**: `kubectl kustomize kustomize/bank-of-anthos/overlays/hackathon/ | kubectl apply -n bank-of-anthos -f -`
  - **Verify**: `kubectl get pods -n bank-of-anthos` shows all pods running

- [ ] Task 5.3: Deploy Online Boutique via kubectl (pre-Argo test)
  - **Depends on**: 5.1
  - **Output**: Online Boutique running in cluster
  - **Commands**: `kubectl kustomize kustomize/online-boutique/overlays/hackathon/ | kubectl apply -n online-boutique -f -`
  - **Verify**: `kubectl get pods -n online-boutique` shows pods running (no paymentservice)

- [ ] Task 5.4: Verify Ingress functionality
  - **Depends on**: 5.2, 5.3
  - **Output**: Working Ingress endpoints
  - **Verify**: 
    - `kubectl get ingress -A` shows both ingresses with IPs
    - Port-forward test: `kubectl port-forward -n bank-of-anthos svc/frontend 8080:80`

- [ ] Task 5.5: Install and configure Argo CD (if not present)
  - **Output**: Argo CD running in cluster
  - **Commands**: Follow Argo CD installation guide
  - **Verify**: `kubectl get pods -n argocd` shows Argo CD components

- [ ] Task 5.6: Deploy Argo CD Applications
  - **Depends on**: 5.5
  - **Output**: Applications managed by Argo CD
  - **Commands**: `kubectl apply -f argocd/applications/`
  - **Verify**: Argo CD UI shows both applications synced

### Phase 5 Checkpoint
- [ ] Service verification: `kubectl get svc -n bank-of-anthos frontend` shows ClusterIP
- [ ] Service verification: `kubectl get svc -n online-boutique | grep -v paymentservice`
- [ ] Ingress verification: Both applications accessible via Ingress endpoints
- [ ] Argo CD verification: Applications show as "Synced" and "Healthy"
- [ ] **Demo ready**: Access both applications through browser via Ingress URLs

## Final Verification
- [ ] All requirements from design doc met:
  - [ ] REQ-001: Bank of Anthos in dedicated namespace
  - [ ] REQ-002: Online Boutique in dedicated namespace
  - [ ] REQ-003: Bank frontend Service is ClusterIP
  - [ ] REQ-004: Boutique frontend-external removed/modified
  - [ ] REQ-005: Ingress resources created
  - [ ] REQ-006: Payment service removed
  - [ ] REQ-007: Kustomize overlays working
  - [ ] REQ-008: Argo CD managing deployments
  - [ ] REQ-009: Frontend deployments running
- [ ] Load generators functioning
- [ ] No payment service pods in Online Boutique
- [ ] Documentation updated with deployment instructions

## Rollback Plan
If issues occur:
1. Via Argo CD: Use UI/CLI to rollback to previous sync
2. Via kubectl: `kubectl delete -f argocd/applications/` then redeploy base manifests
3. Emergency: `kubectl delete namespace bank-of-anthos online-boutique`