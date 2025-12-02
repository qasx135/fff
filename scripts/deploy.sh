#!/bin/bash
set -e

PROJECT_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$PROJECT_ROOT"

echo "🚀 [1/8] Запуск Minikube..."
minikube start --memory=8192 --cpus=4 --driver=docker
minikube addons enable ingress
eval $(minikube docker-env)

echo "📦 [2/8] Добавление Helm-репозиториев (только Prometheus)..."
helm repo add prometheus-community https://prometheus-community.github.io/helm-charts
helm repo update

echo "📦 [3/8] Установка внешних зависимостей..."
helm upgrade --install redis oci://registry-1.docker.io/bitnamicharts/redis --set auth.enabled=false
helm upgrade --install kafka oci://registry-1.docker.io/bitnamicharts/kafka
helm upgrade --install catalog-db oci://registry-1.docker.io/bitnamicharts/postgresql -f infra/postgres/catalog-db-values.yaml
helm upgrade --install watch-db oci://registry-1.docker.io/bitnamicharts/postgresql -f infra/postgres/watch-db-values.yaml
helm upgrade --install mongodb oci://registry-1.docker.io/bitnamicharts/mongodb --set auth.enabled=false --set persistence.size=1Gi --set service.port=27017
echo "📦 [4/8] Установка Graylog через манифест..."
kubectl apply -f infra/graylog.yaml

echo "📈 [5/8] Установка Prometheus + Grafana..."
helm upgrade --install prometheus prometheus-community/kube-prometheus-stack

echo "🔍 [6/8] Установка Jaeger (all-in-one)..."
kubectl create namespace observability --dry-run=client -o yaml | kubectl apply -f -

kubectl apply -f infra/jaeger.yaml

echo "⛵ [7/8] Сборка и установка ваших сервисов..."
docker build -t kostuwan/auth-service:latest ./services/auth-service
docker build -t kostuwan/catalog-service:latest ./services/catalog-service
docker build -t kostuwan/watch-service:latest ./services/watch-service

docker push kostuwan/auth-service:latest
docker push kostuwan/catalog-service:latest
docker push kostuwan/watch-service:latest

helm upgrade --install auth ./helm/auth-service
helm upgrade --install catalog ./helm/catalog-service
helm upgrade --install watch ./helm/watch-service

helm upgrade --install krakend ./helm/krakend

echo
echo "✅ Развёртывание завершено!"
echo "👉 API Gateway: http://anime.local"
echo "👉 Grafana: запустите 'kubectl port-forward -n default svc/prometheus-grafana 3000:80' и откройте http://localhost:3000"
echo "👉 Grafana: запустите 'kubectl port-forward svc/graylog 9000:9000' и откройте http://localhost:9000"
echo "👉 Grafana: запустите 'kubectl port-forward -n observability svc/jaeger 16686:16686' и откройте http://localhost:16686"
