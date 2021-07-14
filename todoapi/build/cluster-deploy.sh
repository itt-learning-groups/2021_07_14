#!/bin/bash

mv ~/.kube/config ~/.kube/work_config

eksctl create cluster \
    --name learnanddevops \
    --version 1.19 \
    --region us-west-2 \
    --nodegroup-name worker-nodes \
    --node-type t3.micro \
    --nodes 3 --nodes-min 1 --nodes-max 4 \
    --managed

aws eks update-kubeconfig --name learnanddevops --region us-west-2
