# Microblogging System

A scalable microblogging platform built with Go, using a microservices architecture. The system consists of four main services: Feed Service, Tweets CRUD, Users CRUD, and Analytics Service, communicating via HTTP and message queues.

## Architecture Overview

The system follows a microservices architecture with the following components:

- **Feed Service**: Handles user timelines and feed generation
- **Tweets CRUD**: Manages tweet creation, reading, updating, and deletion
- **Users CRUD**: Handles user management and relationships
- **Analytics Service**: Processes events and updates caches

### Data Storage

- **PostgreSQL**: Primary data store for tweets and users
- **Redis**: Used for caching timelines, popular tweets, and user data
- **Kafka**: Message queue for event-driven communication between services

## Prerequisites

- Go 1.24.6
- Docker and Docker Compose
- PostgreSQL 13+
- Redis 6+
- Apache Kafka 2.8+
- Minikube (local tests)
- Kubectl

## Getting Started

### 1. Clone the repository

```bash
git clone https://github.com/yourusername/microblogging.git
cd microblogging
```

### 2. Set up configuration

Copy the example configuration file and update it with your settings:

```bash
cp config/config.example.yaml config/config.yaml
# Edit config/config.yaml with your settings
```

### Alternative 1

#### 3. Build services

```bash
docker compose build
```

#### 4. Publish services

```bash
minikube image load microblogging-users:latest
minikube image load microblogging-tweets:latest
minikube image load microblogging-feed:latest
minikube image load microblogging-analytics:latest
minikube image load microblogging-docsify:latest
minikube image load microblogging-swagger-ui:latest
```

#### 5. Deploy services

```bash
kubectl apply -f k8s/
```

### Alternative 2

#### 3. Run deployment script

```bash
./deploy/deploy.sh
```

## Project Structure

```
.
├── cmd/                    # Main applications for the project
│   ├── analytics/          # Analytics service
│   ├── feed/               # Feed service
│   ├── tweets/             # Tweets service
│   └── users/              # Users service
├── pkg/                    # Library code
│   ├── cache/              # Redis cache client
│   ├── config/             # Configuration management
│   ├── database/           # Database clients and models
│   └── queue/              # Kafka message queue client
└── config/                 # Configuration files
    └── *.yaml              # Configuration files
```

## API Documentation

API documentation is available via Swagger UI, Docsify and GitHub Wiki

## Deployment

The application is designed to be deployed on Kubernetes. See the `deploy/` directory for Kubernetes manifests.
