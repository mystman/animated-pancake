#!/bin/sh
eval $(minikube -p minikube docker-env)
docker build -t $1:$2 .
eval $(minikube -p minikube docker-env -u)