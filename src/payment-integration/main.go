package main

import (
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/gke-hackathon/payment-integration/logging"
	"github.com/gke-hackathon/payment-integration/metrics"
	"github.com/gke-hackathon/payment-integration/server"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
)

const defaultPort = "50051"

func main() {
	logger := logging.NewLogger("payment-integration")

	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	logger.Info("Starting payment integration service", map[string]interface{}{"port": port})

	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", port))
	if err != nil {
		logger.Fatal("Failed to listen", err)
	}

	grpcServer := grpc.NewServer()

	paymentServer := server.NewPaymentServer()
	server.RegisterPaymentServiceServer(grpcServer, paymentServer)

	// Register health service
	healthServer := health.NewServer()
	grpc_health_v1.RegisterHealthServer(grpcServer, healthServer)
	healthServer.SetServingStatus("", grpc_health_v1.HealthCheckResponse_SERVING)
	healthServer.SetServingStatus("hipstershop.PaymentService", grpc_health_v1.HealthCheckResponse_SERVING)

	// Register reflection service for grpcurl and other tools
	reflection.Register(grpcServer)

	// Start HTTP server for health checks and metrics
	go startHTTPServer(logger)

	// Handle graceful shutdown
	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
		<-sigChan
		logger.Info("Received shutdown signal, gracefully stopping...", nil)
		grpcServer.GracefulStop()
	}()

	logger.Info("Payment integration service listening", map[string]interface{}{"address": lis.Addr().String()})
	if err := grpcServer.Serve(lis); err != nil {
		logger.Fatal("Failed to serve", err)
	}
}

// startHTTPServer starts the HTTP server for health checks and metrics
func startHTTPServer(logger *logging.Logger) {
	httpPort := os.Getenv("HTTP_PORT")
	if httpPort == "" {
		httpPort = "8080"
	}

	http.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	http.HandleFunc("/readyz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	// Add Prometheus metrics endpoint
	http.Handle("/metrics", metrics.PrometheusHandler())

	logger.Info("Starting HTTP server for health checks", map[string]interface{}{"port": httpPort})
	if err := http.ListenAndServe(":"+httpPort, nil); err != nil {
		logger.Error("HTTP server failed", err, nil)
	}
}
