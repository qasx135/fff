#!/bin/bash
set -e

PROJECT_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$PROJECT_ROOT"

echo "🚀 [1/8] Запуск Minikube..."
minikube start --cpus=4 --driver=docker
minikube addons enable ingress
minikube addons enabe ingress-dns
minikube addons enable metrics-server
minikube addons enable dashboard
minikube addons enable default-storage-class
minikube addons enable storage-provisioner
eval $(minikube docker-env)

echo "📦 [3/8] Установка внешних зависимостей..."
helm upgrade --install redis oci://registry-1.docker.io/bitnamicharts/redis --set auth.enabled=false \
    --set master.persistence.size=1Gi \
    -- set replica.persistence.size=1Gi

kubectl wait --for=condition=ready pod -l app.kubernetes.io/name=redis --timeout=100s

helm upgrade --install mongodb oci://registry-1.docker.io/bitnamicharts/mongodb \
    --set auth.enabled=false \
    --set persistence.size=1Gi \
    --set service.port=27017
kubectl wait --for=condition=ready pod -l app.kubernetes.io/name=mongodb --timeout=100s

helm upgrade --install catalog-db oci://registry-1.docker.io/bitnamicharts/postgresql -f infra/postgres/catalog-db-values.yaml
kubectl wait --for=condition=ready pod -l app.kubernetes.io/instance=catalog-db --timeout=100s
helm upgrade --install watch-db oci://registry-1.docker.io/bitnamicharts/postgresql -f infra/postgres/watch-db-values.yaml
kubectl wait --for=condition=ready pod -l app.kubernetes.io/instance=watch-db --timeout=100s
helm upgrade --install catalog-db-exporter ./helm/postgres-exporter -f ./infra/postgres-exporter/catalog-db-values.yaml
kubectl wait --for=condition=ready pod -l app=catalog-db-exporter --timeout=100s
helm upgrade --install watch-db-exporter ./helm/postgres-exporter -f ./infra/postgres-exporter/watch-db-values.yaml
kubectl wait --for=condition=ready pod -l app=watch-db-exporter --timeout=100s

helm upgrade --install kafka ./helm/kafka
kubectl wait --for=condition=ready pod -l component=kafka --timeout=100s

helm upgrade --install opensearch ./helm/opensearch
kubectl wait --for=condition=ready pod -l app=opensearch --timeout=100s
helm upgrade --install graylog ./helm/graylog
kubectl wait --for=condition=ready pod -l app=graylog --timeout=100s

helm upgrade --install prometheus ./helm/prometheus
kubectl wait --for=condition=ready pod -l app=prometheus --timeout=100s
helm upgrade --install grafana ./helm/grafana
kubectl wait --for=condition=ready pod -l app=grafana --timeout=100s

helm upgrade --install jaeger ./helm/jaeger
kubectl wait --for=condition=ready pod -l app=jaeger --timeout=100s

echo "⛵ [7/8] Сборка и установка ваших сервисов..."
docker build -t kostuwan/auth-service:latest ./services/auth-service
docker build -t kostuwan/catalog-service:latest ./services/catalog-service
docker build -t kostuwan/watch-service:latest ./services/watch-service

docker push kostuwan/auth-service:latest
docker push kostuwan/catalog-service:latest
docker push kostuwan/watch-service:latest

helm upgrade --install auth-service ./helm/auth-service
kubectl wait --for=condition=ready pod -l app=auth-servoce --timeout=100s
helm upgrade --install catalog-service ./helm/catalog-service
kubectl wait --for=condition=ready pod -l app=catalog-service --timeout=100s
helm upgrade --install watch-service ./helm/watch-service
kubectl wait --for=condition=ready pod -l app=watch-service --timeout=100s

helm upgrade --install krakend ./helm/krakend
kubectl wait --for=condition=ready pod -l app=krakend --timeout=100s

echo
echo "✅ Развёртывание завершено!"
