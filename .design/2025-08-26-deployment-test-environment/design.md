# Design: GKE Hackathon Test Environment Deployment

## Layer 1: Problem & Requirements

### Problem Statement
We need to deploy two production-ready demo applications (Bank of Anthos and Online Boutique) for the GKE Turns 10 hackathon, where we'll build AI agents that enhance these microservices using the A2A protocol. The deployment must be customizable to disable certain services (frontend and payment) that will be replaced with enhanced versions or exposed differently.

### Current State
- Two separate Google sample applications exist with independent deployment structures
- Bank of Anthos: 9-service banking application with Python/Java microservices
- Online Boutique: 11-service e-commerce application with polyglot microservices
- Both use frontend services as primary entry points with LoadBalancer/NodePort exposure
- Boutique includes a payment service that processes mock credit card transactions
- No unified deployment strategy or GitOps configuration exists

### Requirements
#### Functional
- REQ-001: The system SHALL deploy Bank of Anthos to a dedicated namespace
- REQ-002: The system SHALL deploy Online Boutique to a dedicated namespace  
- REQ-003: WHEN deploying Bank of Anthos THEN the frontend Service type SHALL be changed from LoadBalancer to ClusterIP
- REQ-004: WHEN deploying Online Boutique THEN the frontend-external Service SHALL be removed or changed to ClusterIP
- REQ-005: The system SHALL add Ingress resources to expose both frontend applications externally
- REQ-006: WHEN deploying Online Boutique THEN the payment service SHALL be removed completely (both Deployment and Service)
- REQ-007: The system SHALL use Kustomize overlays to modify base manifests without altering upstream code
- REQ-008: The system SHALL use Argo CD Applications for GitOps-based deployment
- REQ-009: Both frontend Deployments SHALL continue running with their application logic intact

#### Non-Functional
- Performance: Applications must handle standard load testing scenarios
- Security: Services must maintain JWT authentication (Bank of Anthos) and session management
- Usability: Deployment must be reproducible and declarative via GitOps
- Extensibility: Configuration must support future AI agent integration

### Constraints
- Cannot modify upstream application source code (hackathon requirement)
- Must use GKE as the target platform
- Must support A2A protocol integration for agent communication
- Frontend services must be accessible via HTTP/HTTPS ingress

### Success Criteria
- Both applications deploy successfully to GKE
- Frontend applications are accessible via Ingress endpoints
- Bank of Anthos frontend service uses ClusterIP instead of LoadBalancer
- Online Boutique functions without payment service
- All microservices (including frontends) communicate correctly
- Argo CD successfully manages application lifecycle

## Layer 2: Functional Specification

### User Workflows
1. **Initial Deployment**
   - DevOps pushes Kustomize configurations to Git → Argo CD detects changes → Applications sync to cluster → Services become available via Ingress

2. **Service Access**
   - User accesses Bank of Anthos → Ingress routes to frontend service → Frontend renders UI → User interacts with banking services
   - User accesses Online Boutique → Ingress routes to frontend service → Frontend renders shop UI → Checkout bypasses payment (handled by AI agent)

3. **AI Agent Integration**
   - Agent connects to service APIs → Uses A2A protocol for communication → Enhances functionality without code changes

### External Interfaces
#### Ingress Endpoints
- Bank of Anthos: `bank.hackathon.example.com`
- Online Boutique: `boutique.hackathon.example.com`

#### Service APIs (for AI agents)
- Bank of Anthos REST/gRPC endpoints remain unchanged
- Boutique gRPC services accessible for agent integration
- Payment processing redirected to custom AI payment agent

### Alternatives Considered
| Option | Pros | Cons | Why Not Chosen |
|--------|------|------|----------------|
| Helm Charts | Parameterized deployments | Complex value management, less GitOps-friendly | Kustomize provides cleaner overlays |
| Manual Deployment | Simple initial setup | No GitOps, not reproducible | Violates automation requirements |
| Service Mesh Frontend | Advanced routing capabilities | Overhead for demo apps | Unnecessary complexity |

## Layer 3: Technical Specification

### Architecture
```
┌─────────────────────────────────────────────────────────┐
│                     GKE Cluster                          │
├─────────────────────────────────────────────────────────┤
│  ┌─────────────────────────────────────────────────┐    │
│  │              Ingress Controller                  │    │
│  └────────────┬─────────────────┬──────────────────┘    │
│               │                 │                        │
│  ┌────────────▼──────┐ ┌───────▼──────────┐            │
│  │  bank-of-anthos   │ │ online-boutique   │            │
│  │    namespace      │ │   namespace       │            │
│  ├──────────────────┤ ├──────────────────┤            │
│  │ - userservice    │ │ - cartservice    │            │
│  │ - ledger-writer  │ │ - checkoutservice│            │
│  │ - balance-reader │ │ - currencyservice│            │
│  │ - contacts       │ │ - emailservice   │            │
│  │ - transaction-   │ │ - productcatalog │            │
│  │   history        │ │ - recommendation │            │
│  │ - accounts-db    │ │ - shipping       │            │
│  │ - ledger-db      │ │ - adservice      │            │
│  │ - loadgenerator  │ │ - loadgenerator  │            │
│  │ - frontend       │ │ - frontend       │            │
│  └──────────────────┘ │ ❌ paymentservice│            │
│                        └──────────────────┘            │
├─────────────────────────────────────────────────────────┤
│                    Argo CD                              │
│  ┌──────────────┐  ┌──────────────────┐                │
│  │ Application: │  │   Application:   │                │
│  │ bank-of-     │  │ online-boutique  │                │
│  │   anthos     │  │                  │                │
│  └──────────────┘  └──────────────────┘                │
└─────────────────────────────────────────────────────────┘
```

### Code Change Analysis
| Component | Action | Justification |
|-----------|--------|---------------|
| kustomize/bank-of-anthos/ | Create | New Kustomize structure for Bank of Anthos |
| kustomize/online-boutique/ | Create | New Kustomize structure for Online Boutique |
| argocd/applications/ | Create | Argo CD Application manifests |
| ingress/configurations/ | Create | Ingress resource definitions |
| Bank frontend Service type | Patch | Change from LoadBalancer to ClusterIP |
| Boutique frontend-external Service | Remove/Patch | Remove or change to ClusterIP |
| Boutique payment service | Remove | Will be replaced by AI agent |

### Code to Remove/Modify
- **Frontend Service Types** (both applications)
  - Why it's changing: Direct LoadBalancer exposure not needed
  - What replaces it: ClusterIP service with Ingress controller routing
  - Migration path: Frontend services remain but only cluster-internal

- **Frontend-external Service** (Online Boutique)  
  - Why it's obsolete: Redundant with Ingress routing
  - What replaces it: Single frontend service with ClusterIP
  - Migration path: Remove or convert to ClusterIP

- **Payment Service** (Online Boutique)
  - Why it's obsolete: Will be replaced with AI-enhanced payment agent
  - What replaces it: Custom payment agent using A2A protocol
  - Migration path: Checkout service will handle missing payment service gracefully

### Implementation Approach

#### Components

**Kustomize Base Structure** (`kustomize/`)
- Current role: None (new)
- Planned changes: Create overlay structure
- Integration approach: Reference upstream manifests
- Example logic (pseudocode):
  ```
  for each application:
    create base/kustomization.yaml referencing upstream
    create overlays/hackathon/ with patches
    apply strategic merge patches to remove services
    add ingress resources
  ```

**Bank of Anthos Overlay** (`kustomize/bank-of-anthos/overlays/hackathon/`)
- Components to patch:
  - Patch frontend Service type to ClusterIP
  - Add ingress resource pointing to frontend service
  - Keep all deployments running
  
**Online Boutique Overlay** (`kustomize/online-boutique/overlays/hackathon/`)
- Components to patch:
  - Remove or patch frontend-external Service
  - Remove paymentservice Deployment
  - Remove paymentservice Service
  - Add ingress resource pointing to frontend service
  - Keep frontend deployment running

**Argo CD Applications** (`argocd/applications/`)
- Bank of Anthos Application manifest
- Online Boutique Application manifest
- Configuration:
  ```
  source:
    repoURL: [git-repo]
    path: kustomize/[app]/overlays/hackathon
    targetRevision: HEAD
  destination:
    namespace: [app-namespace]
  syncPolicy:
    automated:
      prune: true
      selfHeal: true
  ```

#### Data Models
No schema changes required - using existing application databases

#### Security
- Maintain JWT tokens for Bank of Anthos authentication
- Preserve session management for Online Boutique
- Ingress TLS termination for HTTPS access
- Network policies to restrict inter-service communication

### Testing Strategy
- Unit tests: Not applicable (configuration only)
- Integration tests: 
  - Verify frontend services are accessible via Ingress
  - Test checkout flow without payment service
  - Validate ingress routing to frontend services
- E2E tests:
  - Full user journey through Bank of Anthos frontend via Ingress
  - Shopping cart and checkout (minus payment) in Online Boutique frontend
  - Load generator functionality validation

### Rollout Plan
1. **Phase 1**: Deploy Kustomize configurations to Git repository
2. **Phase 2**: Install Argo CD on GKE cluster
3. **Phase 3**: Create Argo CD Applications pointing to Kustomize overlays
4. **Phase 4**: Verify deployments and service health
5. **Phase 5**: Configure Ingress DNS and TLS
6. **Rollback**: Argo CD revision rollback or manual manifest application

### Directory Structure
```
gke-hackathon/
├── kustomize/
│   ├── bank-of-anthos/
│   │   ├── base/
│   │   │   └── kustomization.yaml
│   │   └── overlays/
│   │       └── hackathon/
│   │           ├── kustomization.yaml
│   │           ├── patches/
│   │           │   └── frontend-service-patch.yaml
│   │           └── resources/
│   │               └── ingress.yaml
│   └── online-boutique/
│       ├── base/
│       │   └── kustomization.yaml
│       └── overlays/
│           └── hackathon/
│               ├── kustomization.yaml
│               ├── patches/
│               │   ├── frontend-service-patch.yaml
│               │   └── remove-payment.yaml
│               └── resources/
│                   └── ingress.yaml
└── argocd/
    └── applications/
        ├── bank-of-anthos.yaml
        └── online-boutique.yaml
```

### Next Steps
1. Create Kustomize directory structure
2. Write Kustomization files referencing upstream manifests
3. Create strategic merge patches for service removal
4. Define Ingress resources with proper routing rules
5. Create Argo CD Application manifests
6. Deploy and validate in GKE cluster
7. Integrate AI agents using A2A protocol