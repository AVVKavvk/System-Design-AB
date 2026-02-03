#!/bin/bash

echo "🧹 Cleaning up leader election demo..."

# Delete deployment
echo "Deleting leader election app..."
kubectl delete -f k8s/deployment.yaml 2>/dev/null || true

# Delete etcd
echo "Deleting etcd..."
kubectl delete -f k8s/etcd-statefulset.yaml 2>/dev/null || true

# Wait a bit for resources to be deleted
sleep 5

# Check if anything is left
REMAINING=$(kubectl get pods -l app=leader-election,app=etcd 2>/dev/null | wc -l)
if [ "$REMAINING" -gt 1 ]; then
    echo "⚠️  Some pods are still terminating. Waiting..."
    kubectl wait --for=delete pod -l app=leader-election --timeout=60s 2>/dev/null || true
    kubectl wait --for=delete pod -l app=etcd --timeout=60s 2>/dev/null || true
fi

echo "✅ Cleanup complete!"
echo ""
echo "Minikube is still running. To stop it:"
echo "  minikube stop"
echo ""
echo "To delete the entire cluster:"
echo "  minikube delete"