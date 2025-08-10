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
kubectl -n microblogging wait --for=condition=ready pod -l app=postgres --timeout=120s
kubectl -n microblogging wait --for=condition=ready pod -l app=redis --timeout=120s
kubectl -n microblogging wait --for=condition=ready pod -l app=zookeeper --timeout=120s
kubectl -n microblogging wait --for=condition=ready pod -l app=kafka --timeout=120s

echo "Building and pushing service images..."
# docker compose build docsify
# docker compose build users-service
# docker compose build tweets-service
# docker compose build feed-service
# docker compose build analytics-service

# Replace minikube with docker push when not using minikube
# minikube image load microblogging-docsify:latest
# minikube image load microblogging-users-service:latest
# minikube image load microblogging-tweets-service:latest
# minikube image load microblogging-feed-service:latest
# minikube image load microblogging-analytics-service:latest

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
