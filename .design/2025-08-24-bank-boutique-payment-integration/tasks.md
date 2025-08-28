# Implementation Tasks: Bank-Boutique Payment Integration

## Overview
Building a payment integration service that bridges Online Boutique's gRPC payment interface with Bank of Anthos' REST banking API. The service will handle JWT authentication, credit card to bank account mapping, and money format conversion.

The implementation is divided into 5 phases, progressing from a minimal skeleton to full production deployment.

## Phase 1: Minimal Skeleton Service
**Goal**: Create a basic gRPC service that responds to payment requests with hardcoded success
**Demo**: "At standup, I can show: the CheckoutService can call our new payment service and get a successful response"

### Tasks
- [x] Task 1.1: Create project structure and proto definitions
  - **Output**: Basic Go project with proto files
  - **Files**: 
    - `src/payment-integration/go.mod`
    - `src/payment-integration/proto/payment.proto` 
    - `src/payment-integration/Dockerfile`
  - **Verify**: `go mod tidy` runs successfully

- [x] Task 1.2: Implement minimal gRPC server
  - **Depends on**: 1.1
  - **Output**: gRPC server that listens on port 50051
  - **Files**: 
    - `src/payment-integration/main.go`
    - `src/payment-integration/server/server.go`
  - **Verify**: Server starts with `go run main.go`

- [x] Task 1.3: Implement hardcoded Charge RPC
  - **Depends on**: 1.2
  - **Output**: Charge method returns fixed transaction ID
  - **Files**: `src/payment-integration/server/charge.go` (integrated in server.go)
  - **Verify**: Server responds with TXID-12345

- [x] Task 1.4: Create Kubernetes deployment manifest
  - **Output**: Deployable container on K8s
  - **Files**: 
    - `k8s/kustomize/payment-integration/` structure
    - Updated `k8s/kustomize/online-boutique/overlays/hackathon/kustomization.yaml`
  - **Remove**: Removed reference to remove-payment.yaml patch
  - **Verify**: Kustomize structure ready for deployment

### Phase 1 Checkpoint
- [x] Run build: `docker build -t payment-integration src/payment-integration/` ✓
- [x] Run local test: `go test ./...` in src/payment-integration ✓
- [x] Deploy to cluster: Ready for deployment via kustomize
- [x] Manual verification: Docker container runs successfully on port 50051
- [x] **Demo ready**: Payment service responds with hardcoded TXID-12345

## Phase 2: Money Conversion & Basic Mapping
**Goal**: Correctly handle money format conversion and implement credit card to account mapping
**Demo**: "At standup, I can show: payments with different amounts are converted correctly and card numbers map to accounts"

### Tasks
- [x] Task 2.1: Implement money converter utility
  - **Output**: Bidirectional conversion between formats
  - **Files**: 
    - `src/payment-integration/converter/money.go`
    - `src/payment-integration/converter/money_test.go`
  - **Verify**: Unit tests pass for edge cases (0, large numbers) ✓

- [x] Task 2.2: Implement account mapper
  - **Output**: Credit card to bank account mapping logic
  - **Files**: 
    - `src/payment-integration/mapper/account.go`
    - `src/payment-integration/mapper/account_test.go`
  - **Verify**: Last 10 digits correctly extracted ✓

- [x] Task 2.3: Integrate converter and mapper into Charge RPC
  - **Depends on**: 2.1, 2.2
  - **Output**: Real amount processing and account mapping
  - **Files**: Updated `src/payment-integration/server/server.go`
  - **Verify**: Different card numbers produce different account IDs ✓

- [x] Task 2.4: Add environment configuration
  - **Output**: Configurable merchant account details
  - **Files**: 
    - `src/payment-integration/config/config.go`
    - `k8s/kustomize/payment-integration/base/configmap.yaml`
  - **Verify**: Config values readable from environment ✓

### Phase 2 Checkpoint
- [x] Run unit tests: `go test ./converter ./mapper` ✓ All tests pass
- [x] Build and deploy updated service ✓ Via Tilt
- [x] Manual verification: Log output shows correct conversions ✓
- [x] **Demo ready**: Different payment amounts (£1265.24) and card numbers (*6079) processed correctly ✓

## Phase 3: Service-to-Service Authentication
**Goal**: Implement service account authentication for Bank API calls
**Demo**: "At standup, I can show: payment service authenticates with Bank API using service credentials"

### Tasks
- [x] Task 3.1: Implement service token generator
  - **Output**: Generate/load service JWT tokens for Bank API
  - **Files**: 
    - `src/payment-integration/auth/service_auth.go`
    - `src/payment-integration/auth/service_auth_test.go`
  - **Verify**: Can generate valid service tokens

- [x] Task 3.2: Create service credentials configuration
  - **Output**: Service account credentials in Kubernetes secret
  - **Files**: 
    - `k8s/kustomize/payment-integration/base/service-secret.yaml`
    - Update deployment to mount secret
  - **Verify**: Credentials accessible in pod

- [x] Task 3.3: Integrate service auth into Bank client
  - **Depends on**: 3.1, 3.2
  - **Output**: All Bank API calls include service JWT
  - **Files**: Update `src/payment-integration/bank/client.go` (prep for Phase 4)
  - **Verify**: Service token included in API headers

- [x] Task 3.4: Test service authentication flow
  - **Output**: Verify end-to-end auth without user tokens
  - **Files**: Create integration test
  - **Verify**: Boutique can process payments without user auth

### Phase 3 Checkpoint
- [x] Run auth tests: `go test ./auth`
- [x] Verify service token generation
- [x] Test payment flow without user authentication
- [x] **Demo ready**: Show service-to-service authentication working

## Phase 4: Bank API Integration
**Goal**: Connect to real Bank of Anthos API to process actual transfers
**Demo**: "At standup, I can show: a boutique purchase creates a real bank transaction visible in account history"

### Tasks
- [x] Task 4.1: Implement Bank API client
  - **Output**: REST client for Bank transactions endpoint
  - **Files**: 
    - `src/payment-integration/bank/client.go`
    - `src/payment-integration/bank/models.go`
  - **Verify**: Can construct proper transaction requests

- [x] Task 4.2: Add transaction ID generation
  - **Output**: Unique UUIDs for each transaction
  - **Files**: `src/payment-integration/utils/uuid.go`
  - **Verify**: No duplicate IDs generated

- [x] Task 4.3: Integrate Bank client into Charge RPC
  - **Depends on**: 4.1, 4.2
  - **Output**: Real bank transactions on payment
  - **Files**: Update `src/payment-integration/server/charge.go`
  - **Verify**: Transaction appears in Bank logs

- [x] Task 4.4: Implement error handling for Bank failures
  - **Output**: Graceful handling of Bank API errors
  - **Files**: 
    - `src/payment-integration/bank/errors.go`
    - Update charge.go with retry logic
  - **Verify**: Appropriate errors returned to CheckoutService

- [x] Task 4.5: Add Bank service discovery configuration
  - **Output**: Correct routing to Bank services
  - **Files**: Update ConfigMap with Bank endpoints
  - **Verify**: Can reach Bank API from pod

### Phase 4 Checkpoint
- [x] Run integration tests with mock Bank
- [x] Deploy and test with real Bank of Anthos
- [x] Verify transaction in Bank UI
- [x] Check insufficient funds handling
- [x] **Demo ready**: Complete end-to-end purchase with bank debit

## Phase 5: Production Readiness & Cleanup
**Goal**: Remove old payment service and ensure production quality
**Demo**: "At standup, I can show: the old mock payment service is gone and everything still works"

### Tasks
- [x] Task 5.1: Add comprehensive logging
  - **Output**: Transaction tracing across systems
  - **Files**: 
    - `src/payment-integration/logging/logger.go`
    - Update all components with logging
  - **Verify**: Can trace payment flow in logs

- [x] Task 5.2: Add metrics and health checks
  - **Output**: Service observability
  - **Files**: 
    - `src/payment-integration/metrics/metrics.go`
    - Add `/healthz` endpoint
  - **Verify**: Prometheus can scrape metrics

- [x] Task 5.3: Remove old payment service code
  - **Output**: Clean codebase without mock service
  - **Remove**: 
    - `_external/microservices-demo/src/paymentservice/` directory
    - Old paymentservice references in manifests
  - **Verify**: No references to old service remain

- [x] Task 5.4: Add rate limiting
  - **Output**: Protection against payment spam
  - **Files**: `src/payment-integration/middleware/ratelimit.go`
  - **Verify**: Requests throttled after limit

- [x] Task 5.5: Final documentation
  - **Output**: Complete service documentation
  - **Files**: 
    - `src/payment-integration/README.md`
    - API documentation
  - **Verify**: Another developer can understand the service

### Phase 5 Checkpoint
- [x] Run full test suite: `go test ./...`
- [x] Load test the service
- [x] Verify old payment service fully removed
- [x] Check all logging and metrics work
- [x] **Demo ready**: Production-quality payment integration with full observability

## Final Verification
- [ ] All requirements from design doc met:
  - [ ] REQ-001: Accepts gRPC payment requests
  - [ ] REQ-002: Debits Bank of Anthos accounts
  - [ ] REQ-003: Validates JWT tokens
  - [ ] REQ-004: Maps card numbers to accounts
  - [ ] REQ-005: Handles insufficient funds
  - [ ] REQ-006: Converts money formats
- [ ] All obsolete code removed (old paymentservice)
- [ ] Tests comprehensive (unit, integration, e2e)
- [ ] Documentation complete
- [ ] Performance < 500ms p99 latency
- [ ] Zero modification to core Bank/Boutique code