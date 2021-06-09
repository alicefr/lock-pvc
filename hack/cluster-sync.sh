#!/bin/bash

set -x

IMAGE=controller
minikube image rm $IMAGE
kubectl delete deploy lock-pvc-lock-pvc-controller-manager
set -e  -o pipefail
make
make docker-build
minikube image load $IMAGE
make deploy
kubectl config set-context lock-pvc --namespace=lock-pvc-system --cluster=minikube --user=minikube
kubectl config use-context lock-pvc

