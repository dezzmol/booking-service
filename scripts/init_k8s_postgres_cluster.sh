#!/bin/bash

#echo "Останавливаем Deployment..."
#kubectl delete -n booking-service-db deployment postgres
#
#echo "Удаляем данные..."
#kubectl delete -n booking-service-db pvc payment-postgres-pvc
#
#echo "Удаляем Service..."
#kubectl delete -n booking-service-db service payment-postgres
#
#echo "Удаляем Secret"
#kubectl delete -n booking-service-db secret payment-postgres-secret

echo "Ожидаем удаления..."
sleep 5

echo "Запускаем заново..."
kubectl apply -f ./.k8s/postgres/deployment.yaml
kubectl apply -f ./.k8s/postgres/service.yaml
kubectl apply -f ./.k8s/postgres/secret.yaml
kubectl apply -f ./.k8s/postgres/pvc.yaml

echo "Ждем запуска..."
kubectl wait --for=condition=ready pod -l app=payment-postgres --timeout=60s

echo "Готово!"