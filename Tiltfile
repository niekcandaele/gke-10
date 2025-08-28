# Tiltfile for GKE Hackathon Development
# This file configures Tilt for rapid development of the payment-integration service

# Allow the GKE cluster context
allow_k8s_contexts('gke_takaro-dev_europe-west1_takaro-dev')

# Configure default registry to use ttl.sh (ephemeral registry)
# This registry allows anonymous pushes and images expire after 24 hours
default_registry('ttl.sh/gke-hackathon-' + str(local('echo $USER', quiet=True)).strip())

# Load configuration (if exists)
config = {}
config_path = './tilt_config.json'
if os.path.exists(config_path):
    config = read_json(config_path)

# Payment Integration Service
print('🚀 Building payment-integration service...')

# Build the payment-integration image
# Tilt will automatically rebuild when files change
docker_build(
    'payment-integration',
    './src/payment-integration',
    dockerfile='./src/payment-integration/Dockerfile'
)

# Deploy payment-integration using kustomize
k8s_yaml(kustomize('./k8s/kustomize/payment-integration/overlays/hackathon'))

# Configure the payment-integration resource
k8s_resource(
    'payment-integration',
    port_forwards='50051:50051',  # Forward gRPC port
    labels=['payment'],
    resource_deps=[]
)

# Optional: Deploy the full Online Boutique stack (commented out by default)
# Uncomment the following lines if you want to deploy the entire application
if config.get('deploy_full_stack', False):
    print('📦 Deploying full Online Boutique stack...')
    k8s_yaml(kustomize('./k8s/kustomize/online-boutique/overlays/hackathon'))
    
    # Port forwards for key services
    k8s_resource('frontend', port_forwards='8080:8080', labels=['boutique'])
    k8s_resource('checkoutservice', labels=['boutique'])
    
    print('🏦 Deploying Bank of Anthos...')
    k8s_yaml(kustomize('./k8s/kustomize/bank-of-anthos/overlays/hackathon'))
    k8s_resource('frontend', port_forwards='8081:8080', labels=['bank'], workload='frontend', new_name='bank-frontend')

print('✅ Tilt configuration loaded!')
print('🎯 Access Tilt UI at http://localhost:10350')
print('🔌 Payment service will be available at localhost:50051')