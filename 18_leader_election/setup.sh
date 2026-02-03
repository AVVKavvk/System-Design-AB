#!/bin/bash

set -e

echo "🚀 Leader Election Demo - Setup Script"
echo "========================================"
echo ""

# Colors for output
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Check if minikube is running
if ! minikube status &> /dev/null; then
    echo -e "${YELLOW}Starting Minikube...${NC}"
    minikube start --cpus=4 --memory=4096
else
    echo -e "${GREEN}✓ Minikube is already running${NC}"
fi

echo ""
echo -e "${BLUE}Step 1: Deploying etcd...${NC}"
kubectl apply -f k8s/etcd-statefulset.yaml
echo "Waiting for etcd to be ready..."
kubectl wait --for=condition=ready pod/etcd-0 --timeout=120s
echo -e "${GREEN}✓ etcd is ready${NC}"

echo ""
echo -e "${BLUE}Step 2: Building application image...${NC}"
eval $(minikube docker-env)
docker build -t leader-election:latest .
echo -e "${GREEN}✓ Image built${NC}"

echo ""
echo -e "${BLUE}Step 3: Deploying leader election app (3 replicas)...${NC}"
kubectl apply -f k8s/deployment.yaml
echo "Waiting for pods to be ready..."
kubectl wait --for=condition=ready pod -l app=leader-election --timeout=120s
echo -e "${GREEN}✓ All pods are ready${NC}"

echo ""
echo -e "${GREEN}========================================${NC}"
echo -e "${GREEN}🎉 Demo is ready!${NC}"
echo -e "${GREEN}========================================${NC}"
echo ""
echo "Watch the leader election in action:"
echo ""
echo -e "${YELLOW}Option 1 - Watch all logs together:${NC}"
echo "  kubectl logs -f -l app=leader-election --all-containers=true --prefix=true"
echo ""
echo -e "${YELLOW}Option 2 - Watch individual pods (open 3 terminals):${NC}"
echo "  Terminal 1: kubectl logs -f \$(kubectl get pods -l app=leader-election -o jsonpath='{.items[0].metadata.name}')"
echo "  Terminal 2: kubectl logs -f \$(kubectl get pods -l app=leader-election -o jsonpath='{.items[1].metadata.name}')"
echo "  Terminal 3: kubectl logs -f \$(kubectl get pods -l app=leader-election -o jsonpath='{.items[2].metadata.name}')"
echo ""
echo -e "${YELLOW}Option 3 - Use stern (if installed):${NC}"
echo "  stern leader-election"
echo ""
echo -e "${BLUE}Get pod names:${NC}"
kubectl get pods -l app=leader-election
echo ""
echo -e "${BLUE}Test failover by deleting the leader pod:${NC}"
echo "  kubectl delete pod <pod-name>"
echo ""
echo "To cleanup: ./cleanup.sh"