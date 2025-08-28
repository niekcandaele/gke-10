# Payment Integration Service

## Overview

The Payment Integration Service is a gRPC-based microservice that bridges Online Boutique's checkout system with Bank of Anthos's transaction processing. It replaces the mock payment service with real bank transfers.

## Features

- **gRPC Payment API**: Implements the `PaymentService` interface for seamless integration
- **Bank Integration**: Processes real transactions via Bank of Anthos REST API
- **JWT Authentication**: Service-to-service authentication with account-specific tokens
- **Card-to-Account Mapping**: Maps credit card numbers to bank account numbers
- **Structured Logging**: Comprehensive transaction logging with correlation IDs
- **Rate Limiting**: Per-account rate limiting to prevent abuse
- **Health Checks**: HTTP endpoints for Kubernetes probes
- **Metrics**: Service metrics for monitoring

## Architecture

```
Online Boutique                     Payment Integration              Bank of Anthos
   Frontend                              Service                      Ledger Writer
      │                                     │                              │
      ├─[1. Checkout]─────────►             │                              │
      │                                     │                              │
   Checkout                                 │                              │
   Service                                  │                              │
      │                                     │                              │
      ├─[2. ChargeRequest]────►             │                              │
      │    (gRPC)                           │                              │
      │                                ┌────▼────┐                        │
      │                                │ Validate│                        │
      │                                │ Request │                        │
      │                                └────┬────┘                        │
      │                                     │                              │
      │                                ┌────▼────┐                        │
      │                                │  Map    │                        │
      │                                │Card→Acct│                        │
      │                                └────┬────┘                        │
      │                                     │                              │
      │                                ┌────▼────┐                        │
      │                                │Generate │                        │
      │                                │   JWT   │                        │
      │                                └────┬────┘                        │
      │                                     │                              │
      │                                     ├─[3. POST /transactions]─────►│
      │                                     │     (REST + JWT)             │
      │                                     │                              │
      │                                     │◄─[4. Transaction Response]──┤
      │                                     │                              │
      │◄─[5. ChargeResponse]─────────       │                              │
      │    (Transaction ID)                 │                              │
```

## Configuration

### Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `PORT` | gRPC server port | `50051` |
| `HTTP_PORT` | HTTP server port for health checks | `8080` |
| `MERCHANT_ACCOUNT` | Merchant bank account number | `9999999999` |
| `ROUTING_NUMBER` | Bank routing number | `883745000` |
| `BANK_API_URL` | Bank of Anthos API endpoint | `http://ledgerwriter.bank-of-anthos:8080` |
| `PRIV_KEY_PATH` | Path to JWT private key | `/tmp/.ssh/privatekey` |
| `PUB_KEY_PATH` | Path to JWT public key | `/tmp/.ssh/publickey` |
| `TOKEN_EXPIRY_SECONDS` | JWT token expiry time | `3600` |
| `RATE_LIMIT_PER_MINUTE` | Max transactions per account per minute | `10` |
| `LOG_LEVEL` | Logging level (DEBUG/INFO) | `INFO` |

## API Reference

### gRPC API

#### ChargeRequest
```protobuf
message ChargeRequest {
  Money amount = 1;
  CreditCardInfo credit_card = 2;
}
```

#### ChargeResponse
```protobuf
message ChargeResponse {
  string transaction_id = 1;
}
```

### HTTP Endpoints

- `GET /healthz` - Health check endpoint (returns 200 OK)
- `GET /readyz` - Readiness check endpoint (returns 200 OK)
- `GET /metrics` - Prometheus metrics endpoint

## Card Number Mapping

The service maps credit card numbers to bank accounts using the last 10 digits:

```
Card: 4532 0110 1122 6111
         └─────────┘
      Maps to account: 1011226111
```

## Rate Limiting

Default: 10 transactions per account per minute
- Returns `ResourceExhausted` error when limit exceeded
- Resets every minute
- Configurable via `RATE_LIMIT_PER_MINUTE`

## Logging

Structured JSON logging includes:
- Transaction IDs for tracing
- Request/response times
- Account numbers (masked)
- Error details
- Bank API call duration

Example log:
```json
{
  "timestamp": "2024-08-27T19:33:12Z",
  "level": "INFO",
  "service": "payment-integration",
  "transaction_id": "e736eea3-691f-4889-8449-36528237cdc6",
  "from_account": "1011226111",
  "to_account": "9999999999",
  "amount": 24796,
  "currency": "USD",
  "message": "Bank transaction successful"
}
```

## Development

### Building

```bash
cd src/payment-integration
docker build -t payment-integration:latest .
```

### Testing

```bash
go test ./...
```

### Running Locally

```bash
go run main.go
```

## Deployment

### Kubernetes

Apply the manifests:
```bash
kubectl apply -k k8s/kustomize/payment-integration/overlays/hackathon/
```

### Tilt

The service is integrated with Tilt for local development:
```bash
tilt up
```

## Troubleshooting

### Common Issues

1. **"sender not authenticated" error**
   - Ensure JWT keys are mounted correctly
   - Verify account number matches JWT claims

2. **"insufficient funds" error**
   - Check account balance in Bank of Anthos
   - Verify routing numbers match

3. **Rate limit errors**
   - Wait 1 minute for limit to reset
   - Adjust `RATE_LIMIT_PER_MINUTE` if needed

4. **Bank API connection errors**
   - Verify `BANK_API_URL` is correct
   - Check network policies allow connection
   - Ensure Bank of Anthos is running

### Debug Mode

Set `LOG_LEVEL=DEBUG` for verbose logging including:
- All request/response details
- JWT token generation
- Bank API calls

## Security

- JWT RS256 signatures for authentication
- Account-specific tokens prevent unauthorized transfers
- Rate limiting prevents abuse
- No sensitive data in logs (card numbers masked)

## Monitoring

### Key Metrics

- `total_requests` - Total payment requests
- `successful_requests` - Successful payments
- `failed_requests` - Failed payments
- `avg_latency_ms` - Average response time
- `rejected_requests` - Rate-limited requests

### Health Checks

Configure Kubernetes probes:
```yaml
livenessProbe:
  httpGet:
    path: /healthz
    port: 8080
readinessProbe:
  httpGet:
    path: /readyz
    port: 8080
```

## Contributing

1. Follow existing code patterns
2. Add tests for new features
3. Update documentation
4. Run linting and tests before committing

## License

Part of the GKE Hackathon project