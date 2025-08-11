#!/bin/bash
set -e

# Get the absolute path of the project root
PROJECT_ROOT=$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)

# Apply all Kubernetes manifests
echo "Creating microblogging namespace..."
kubectl apply -f $PROJECT_ROOT/k8s/namespace.yaml

# Create dependencies
echo "Deploying PostgreSQL..."
kubectl apply -f $PROJECT_ROOT/k8s/database-service.yaml

echo "Deploying Redis..."
kubectl apply -f $PROJECT_ROOT/k8s/cache-service.yaml

echo "Deploying Kafka and Zookeeper..."
kubectl apply -f $PROJECT_ROOT/k8s/queue-service.yaml

echo "Waiting for infrastructure to be ready..."
kubectl -n microblogging wait --for=condition=ready pod -l app=postgres-primary --timeout=120s
kubectl -n microblogging wait --for=condition=ready pod -l app=postgres-replica --timeout=120s
kubectl -n microblogging wait --for=condition=ready pod -l app=redis --timeout=120s
kubectl -n microblogging wait --for=condition=ready pod -l app=zookeeper --timeout=120s
kubectl -n microblogging wait --for=condition=ready pod -l app=kafka --timeout=120s

# Function to check if a service should be built
should_build_service() {
    local service=$1
    # If --skip-build-and-publish is not specified, build all
    if [[ ! " ${@:1} " =~ " --skip-build-and-publish " ]]; then
        return 0  # build this service
    fi
    # If no services specified after --skip-build-and-publish, skip all
    if [[ $# -le 1 || " ${@:1} " =~ " --skip-build-and-publish $ " ]]; then
        return 1  # skip this service
    fi
    # Check if this service is in the skip list
    if [[ " ${@:1} " =~ " --skip-build-and-publish " ]]; then
        local skip_services="${@#*--skip-build-and-publish}"
        if [[ " $skip_services " =~ " $service " ]]; then
            return 1  # skip this service
        fi
    fi
    return 0  # build this service
}

# Build and push services if not skipped
echo "Building and pushing service images..."

if should_build_service "docsify" "$@"; then
    docker compose build docsify --no-cache
    minikube image load microblogging-docsify:latest
else
    echo "Skipping docsify build as requested..."
fi

if should_build_service "users-service" "$@"; then
    docker compose build users-service --no-cache
    minikube image load microblogging-users-service:latest
else
    echo "Skipping users-service build as requested..."
fi

if should_build_service "tweets-service" "$@"; then
    docker compose build tweets-service --no-cache
    minikube image load microblogging-tweets-service:latest
else
    echo "Skipping tweets-service build as requested..."
fi

if should_build_service "feed-service" "$@"; then
    docker compose build feed-service --no-cache
    minikube image load microblogging-feed-service:latest
else
    echo "Skipping feed-service build as requested..."
fi

if should_build_service "analytics-service" "$@"; then
    docker compose build analytics-service --no-cache
    minikube image load microblogging-analytics-service:latest
else
    echo "Skipping analytics-service build as requested..."
fi

echo "Deploying services..."
kubectl apply -f $PROJECT_ROOT/k8s/feed-service.yaml
kubectl apply -f $PROJECT_ROOT/k8s/tweets-service.yaml
kubectl apply -f $PROJECT_ROOT/k8s/users-service.yaml
kubectl apply -f $PROJECT_ROOT/k8s/analytics-service.yaml

echo "Deploying documentation services..."
kubectl create configmap swagger-files \
  --from-file=feed-service.yaml=./docs/swagger/feed-service.yaml \
  --from-file=tweets-service.yaml=./docs/swagger/tweets-service.yaml \
  --from-file=users-service.yaml=./docs/swagger/users-service.yaml \
  --from-file=analytics-service.yaml=./docs/swagger/analytics-service.yaml \
  --dry-run=client -n microblogging -o yaml > $PROJECT_ROOT/k8s/swagger-configmap.yaml
kubectl apply -f $PROJECT_ROOT/k8s/swagger-configmap.yaml
kubectl apply -f $PROJECT_ROOT/k8s/docs-service.yaml
rm $PROJECT_ROOT/k8s/swagger-configmap.yaml

# echo "Waiting for services to be ready..."
# kubectl -n microblogging wait --for=condition=ready pod -l app=feed-service --timeout=120s
# kubectl -n microblogging wait --for=condition=ready pod -l app=tweets-service --timeout=120s
# kubectl -n microblogging wait --for=condition=ready pod -l app=users-service --timeout=120s
# kubectl -n microblogging wait --for=condition=ready pod -l app=analytics-service --timeout=120s

# echo "Waiting for documentation services to be ready..."
# kubectl -n microblogging wait --for=condition=ready pod -l app=docsify --timeout=120s
# kubectl -n microblogging wait --for=condition=ready pod -l app=swagger-ui --timeout=120s

echo "Setting up Ingress..."
kubectl apply -f $PROJECT_ROOT/k8s/ingress.yaml

echo "Deployment complete!"

# Remove if not using minikube
echo "Setting up minikube tunnel..."
echo "Remember to add host to /etc/hosts"
echo "microblogging.local"
minikube tunnel
