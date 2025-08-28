FROM golang:1.23-alpine AS builder
RUN apk add --no-cache ca-certificates git
RUN apk add build-base

WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o /payment-integration .

FROM alpine:3.18
RUN apk add --no-cache ca-certificates
WORKDIR /app
COPY --from=builder /payment-integration /app/payment-integration
EXPOSE 50051
ENTRYPOINT ["/app/payment-integration"]