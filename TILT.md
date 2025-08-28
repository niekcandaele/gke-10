# Tilt Development Setup

This repository is configured to use [Tilt](https://tilt.dev) for rapid Kubernetes development.

## Prerequisites

1. **Tilt installed**: Already installed at `/home/linuxbrew/.linuxbrew/bin/tilt`
2. **kubectl configured**: Connected to GKE cluster `gke_takaro-dev_europe-west1_takaro-dev`
3. **Docker**: For building images locally

## Quick Start

1. Start Tilt:
   ```bash
   tilt up
   ```

2. Open the Tilt UI in your browser:
   - Automatically opens at http://localhost:10350
   - Or press `space` in the terminal

3. The payment-integration service will be:
   - Built automatically
   - Deployed to the `online-boutique` namespace
   - Available at `localhost:50051` (port forwarded)

## Development Workflow

1. **Make code changes** in `src/payment-integration/`
2. **Tilt automatically**:
   - Detects changes
   - Rebuilds the Docker image
   - Updates the deployment in Kubernetes
   - Shows logs in the UI

3. **Test your changes**:
   - Use `grpcurl` to test: `grpcurl -plaintext localhost:50051 list`
   - Check logs in Tilt UI
   - Monitor pod status

## Configuration

- **Tiltfile**: Main configuration at repository root
- **tilt_config.json**: Optional local overrides (gitignored)
  - Copy `tilt_config.json.example` to get started
  - Set `deploy_full_stack: true` to deploy all services

## Useful Commands

- `tilt up` - Start Tilt and deploy services
- `tilt down` - Stop Tilt and remove deployments
- `tilt ci` - Run in CI mode (exits after deploy)
- `tilt logs -f payment-integration` - Follow logs for a service

## Architecture

Tilt uses an ephemeral registry for fast image builds:
- Images are built locally
- Pushed to an in-cluster registry
- No GCR configuration needed
- Automatically cleaned up on `tilt down`

## Troubleshooting

- **Port already in use**: Kill existing Tilt: `tilt down` or `killall tilt`
- **Build failures**: Check Dockerfile and go.mod dependencies
- **Deploy failures**: Verify kustomize configuration
- **Can't connect**: Ensure kubectl context is correct