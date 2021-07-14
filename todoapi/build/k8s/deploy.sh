#!/bin/bash

kubectl apply -f ./deployment.yaml -n default
kubectl apply -f ./service.yaml -n default
