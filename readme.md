# Инструкция

Поднять кластер БД

```shell
kubectl apply -f ./.k8s/postgres/pvc.yaml
kubectl apply -f ./.k8s/postgres/secret.yaml
kubectl apply -f ./.k8s/postgres/service.yaml
kubectl apply -f ./.k8s/postgres/deployment.yaml
```

Поднять поды

```shell
kubectl apply -f ./.k8s/service.yaml
kubectl apply -f ./.k8s/deployment.yaml
```

Открыть порт

```shell
kubectl port-forward -n booking-service svc/postgres 5434:5434
```

Перезапустить под

```shell
kubectl rollout restart deployment booking-payment
```