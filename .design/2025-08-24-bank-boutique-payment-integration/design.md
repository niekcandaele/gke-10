# Design: Bank-Boutique Payment Integration Service

## Layer 1: Problem & Requirements

### Problem Statement
The GKE hackathon requires connecting two separate microservice applications - Bank of Anthos (banking) and Online Boutique (e-commerce) - to enable purchases at the boutique to deduct funds directly from bank accounts. Currently, these systems operate independently with incompatible payment mechanisms: the boutique uses mock credit card processing while the bank handles account-to-account transfers.

### Current State
- **Online Boutique**: Uses a mock PaymentService that validates credit cards but doesn't process real payments
- **Bank of Anthos**: Processes real bank transfers between accounts with JWT authentication
- **Integration Gap**: No mechanism exists to bridge credit card payments to bank account debits
- **Authentication Mismatch**: Boutique has no authentication while Bank requires JWT tokens
- **Protocol Difference**: Boutique uses gRPC, Bank uses REST/HTTP

### Requirements

#### Functional
- REQ-001: The system SHALL accept payment requests from Online Boutique's CheckoutService via gRPC
- REQ-002: WHEN a customer makes a purchase THEN the payment SHALL be deducted from their Bank of Anthos account
- REQ-003: The system SHALL use service account authentication when calling Bank of Anthos APIs
- REQ-004: The system SHALL map credit card numbers to bank account numbers for payment processing
- REQ-005: WHEN payment fails due to insufficient funds THEN the system SHALL return appropriate error to CheckoutService
- REQ-006: The system SHALL convert between Boutique's Money format and Bank's cents format

#### Non-Functional
- **Performance**: Payment processing latency < 500ms (p99)
- **Security**: All bank transactions must be authenticated with valid JWT tokens
- **Compatibility**: Must not modify core application code (per hackathon rules)
- **Availability**: 99.9% uptime with graceful degradation
- **Observability**: Full transaction tracing across both systems

### Constraints
- Cannot modify existing Bank of Anthos or Online Boutique core services
- Must maintain backward compatibility with existing payment flows
- Must use containerized deployment on GKE
- Should leverage Google AI services for intelligent features

### Success Criteria
- Successful end-to-end purchase flow from boutique to bank debit
- Zero failed transactions due to integration issues
- Complete transaction history visible in both systems
- Clean removal of mock payment service

## Layer 2: Functional Specification

### User Workflows

1. **Purchase with Bank Account**
   - User browses Online Boutique and adds items to cart
   - User proceeds to checkout with "credit card" (mapped to bank account)
   - System authenticates user via Bank of Anthos JWT
   - Payment Integration Service debits bank account
   - User receives order confirmation with bank transaction ID
   - Transaction appears in bank account history

2. **Service Authentication Flow**
   - Payment Integration Service configured with service account credentials
   - Service generates/uses pre-configured JWT token for Bank API calls
   - No authentication required from Online Boutique users
   - Service acts as trusted intermediary between Boutique and Bank
   - All Bank API calls authenticated with service token

3. **Insufficient Funds Handling**
   - Payment request received from CheckoutService
   - Balance check performed via Bank's BalanceReader
   - If insufficient: error returned to CheckoutService
   - User notified of payment failure
   - Cart preserved for retry

### External Interfaces

**gRPC Interface (Boutique-facing)**:
```protobuf
service PaymentService {
    rpc Charge(ChargeRequest) returns (ChargeResponse) {}
}

message ChargeRequest {
    Money amount = 1;
    CreditCardInfo credit_card = 2;
}

message ChargeResponse {
    string transaction_id = 1;
}
```

**REST Interface (Bank-facing)**:
```json
POST /transactions
Authorization: Bearer {jwt_token}
{
    "fromAccountNum": "1234567890",
    "fromRoutingNum": "123456789",
    "toAccountNum": "MERCHANT-001",
    "toRoutingNum": "123456789",
    "amount": 15000,
    "uuid": "unique-transaction-id"
}
```

### Alternatives Considered

| Option | Pros | Cons | Why Not Chosen |
|--------|------|------|----------------|
| Direct Database Integration | Fast, simple | Violates service boundaries, breaks encapsulation | Would require modifying core services |
| Webhook-based Async | Scalable, decoupled | Complex error handling, eventual consistency | Added complexity without clear benefit |
| Shared Message Queue | Reliable, async | Requires new infrastructure, complex retry logic | Over-engineered for current scale |
| API Gateway Pattern | Centralized, flexible | Additional component, potential bottleneck | Adds unnecessary complexity |

## Layer 3: Technical Specification

### Architecture

```
┌─────────────────┐     gRPC      ┌──────────────────────┐     REST/JWT    ┌─────────────────┐
│ CheckoutService │──────────────►│ Payment Integration  │───────────────►│ LedgerWriter    │
│   (Boutique)    │◄──────────────│      Service         │◄───────────────│   (Bank)        │
└─────────────────┘   ChargeResp  └──────────────────────┘   Transaction   └─────────────────┘
```

### Code Change Analysis

| Component | Action | Justification |
|-----------|--------|---------------|
| payment-integration-service | Create | New service to bridge boutique and bank systems |
| src/paymentservice (boutique) | Remove | Replaced by integration service |
| kubernetes-manifests/paymentservice.yaml | Extend | Point to new integration service image |
| config.yaml (bank) | Extend | Add merchant account configuration |
| jwt-secret.yaml | Extend | Share JWT keys with integration service |

### Code to Remove

- **src/paymentservice/** (_external/microservices-demo/src/paymentservice/)
  - Why it's obsolete: Mock service with no real payment processing
  - What replaces it: Payment Integration Service with real bank transfers
  - Migration path: Update Kubernetes manifest to use new service image

### Implementation Approach

#### Components

**Payment Integration Service** (src/payment-integration/)
- Current role: None (new service)
- Planned changes: Bridge between boutique gRPC and bank REST
- Integration approach: Implements PaymentService gRPC interface, calls Bank REST API
- Core logic flow:
  ```
  on ChargeRequest:
    # No JWT required from Boutique
    map credit_card_number to account_id
    convert boutique_money to bank_cents
    generate/retrieve service_jwt_token
    call bank /transactions endpoint with service_jwt
    return transaction_id as ChargeResponse
  ```

**Service Authenticator** (src/payment-integration/auth/)
- Generates or loads service account JWT tokens
- Uses pre-configured credentials for Bank API authentication
- Manages token refresh/expiry
- Acts as service-to-service authentication layer
- No user context required from Boutique

**Money Converter** (src/payment-integration/converter/)
- Converts between formats:
  ```
  boutique_money_to_cents:
    cents = money.units * 100 + money.nanos / 10000000
    
  cents_to_boutique_money:
    units = cents / 100
    nanos = (cents % 100) * 10000000
  ```

**Account Mapper** (src/payment-integration/mapper/)
- Maps credit card numbers to bank accounts
- Initial implementation: last 10 digits as account number
- Future: database-backed mapping table
- Fallback to default merchant account

#### Data Models

**Transaction Mapping**:
```yaml
TransactionMap:
  boutique_order_id: string
  bank_transaction_id: string
  amount_cents: integer
  timestamp: datetime
  status: enum [pending, completed, failed]
```

**Account Mapping** (future enhancement):
```yaml
AccountMap:
  card_number_hash: string
  bank_account_num: string
  bank_routing_num: string
  user_id: string
  created_at: datetime
```

#### Security
- Service-to-service authentication using pre-configured JWT
- Payment service acts as trusted intermediary
- TLS for all service communication
- Secret management via Kubernetes secrets
- No logging of sensitive data (card numbers, account numbers)
- Rate limiting on payment endpoints

### Testing Strategy

**Unit Tests**:
- Service token generation and refresh
- Money conversion accuracy (edge cases: 0, negative, overflow)
- Account mapping logic
- Error handling for bank API failures

**Integration Tests**:
- End-to-end payment flow with mock bank service
- gRPC interface compatibility with CheckoutService
- REST client testing with Bank of Anthos

**E2E Tests**:
- Complete purchase flow from boutique to bank debit
- Insufficient funds scenario
- Authentication failure handling
- Concurrent payment processing

### Rollout Plan

**Phase 1: Foundation (Week 1)**
- Deploy Payment Integration Service with basic bridge functionality
- Implement JWT validation and money conversion
- Simple card-to-account mapping

**Phase 2: Integration (Week 2)**
- Replace boutique PaymentService with integration service
- Configure service discovery and networking
- End-to-end testing with both systems