#!/bin/bash

echo "Удаляем Service"
kubectl delete -n booking-service service booking-service

echo "Удаляем Deployment"
kubectl delete -n booking-service deployment booking-service

echo "Ожидаем удаления..."
sleep 5

echo "Запускаем заново..."
kubectl apply -f ./.k8s/deployment.yaml
kubectl apply -f ./.k8s/service.yaml

echo "Ждем запуска..."
kubectl wait --for=condition=ready pod -l -n booking-service app=booking-service --timeout=60s

echo "Готово"