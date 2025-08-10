# Kubernetes Deployment

## Usefull commands

### Applying manifests

```bash
kubectl apply -f k8s/
```

### Scaling

To scale a service:

```bash
kubectl scale deployment/<service-name> --replicas=3 -n microblogging
```

## Troubleshooting

Check pod status:
```bash
kubectl get pods -n microblogging
```

View logs for a pod:
```bash
kubectl logs -f <pod-name> -n microblogging
```

Describe a pod for more details:
```bash
kubectl describe pod <pod-name> -n microblogging
```

## Cleanup

To delete all resources:

```bash
kubectl delete namespace microblogging
```
