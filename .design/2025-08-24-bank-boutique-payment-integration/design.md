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
- REQ-003: The system SHALL validate Bank of Anthos JWT tokens for authentication
- REQ-004: The system SHALL map credit card numbers to bank account numbers for payment processing
- REQ-005: WHEN payment fails due to insufficient funds THEN the system SHALL return appropriate error to CheckoutService
- REQ-006: The system SHALL convert between Boutique's Money format and Bank's cents format
- REQ-007: The system SHALL support the A2A (Agent-to-Agent) protocol for intelligent payment routing

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
- Must support A2A protocol for agent communication

### Success Criteria
- Successful end-to-end purchase flow from boutique to bank debit
- Zero failed transactions due to integration issues
- Complete transaction history visible in both systems
- Demonstration of AI-powered features (fraud detection, smart routing)
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

2. **Authentication Flow**
   - User logs into Bank of Anthos first
   - JWT token stored in session/cookie
   - Token forwarded to Payment Integration Service
   - Service validates token with Bank's public key
   - Account context extracted for payment processing

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

**A2A Protocol Interface**:
```yaml
agent_communication:
  protocol: a2a
  messages:
    - payment_request
    - fraud_check
    - balance_verification
    - transaction_confirmation
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
                                            │
                                            │ A2A Protocol
                                            ▼
                                   ┌──────────────────┐
                                   │   AI Agents      │
                                   │ - Fraud Detection│
                                   │ - Smart Routing  │
                                   └──────────────────┘
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
    extract jwt_token from context
    validate token with bank public key
    map credit_card_number to account_id
    convert boutique_money to bank_cents
    
    if ai_enabled:
      check fraud_score via a2a
      if fraud_detected: reject
    
    call bank /transactions endpoint
    return transaction_id as ChargeResponse
  ```

**JWT Validator** (src/payment-integration/auth/)
- Validates Bank of Anthos JWT tokens
- Uses RS256 with Bank's public key
- Extracts account context from token
- Pattern follows: _external/bank-of-anthos/src/frontend/frontend.py:616-637

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

**A2A Agent Client** (src/payment-integration/agents/)
- Implements A2A protocol client
- Communicates with AI agents for:
  ```
  fraud_detection:
    send payment_context to agent
    receive risk_score
    if score > threshold: reject
    
  smart_routing:
    send transaction_details
    receive optimal_processing_path
    execute recommended_flow
  ```

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
- JWT validation using Bank of Anthos public key (pattern from userservice.py)
- TLS for all service communication
- Secret management via Kubernetes secrets
- No logging of sensitive data (card numbers, account numbers)
- Rate limiting on payment endpoints

### Testing Strategy

**Unit Tests**:
- JWT validation with valid/invalid/expired tokens
- Money conversion accuracy (edge cases: 0, negative, overflow)
- Account mapping logic
- Error handling for bank API failures

**Integration Tests**:
- End-to-end payment flow with mock bank service
- gRPC interface compatibility with CheckoutService
- REST client testing with Bank of Anthos
- A2A protocol message exchange

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

**Phase 3: AI Enhancement (Week 3)**
- Add A2A protocol support
- Implement fraud detection agent
- Deploy smart routing capabilities

**Feature Flags**:
```yaml
features:
  use_real_bank: true
  enable_fraud_detection: false
  enable_smart_routing: false
  fallback_to_mock: true
```

**Rollback Strategy**:
- Keep original PaymentService deployment yaml
- Feature flag to route to mock service
- Circuit breaker for bank API failures
- Transaction reconciliation job for inconsistencies

**Monitoring**:
- Payment success/failure rates
- Latency percentiles (p50, p95, p99)
- Bank API availability
- JWT validation failures
- A2A agent response times