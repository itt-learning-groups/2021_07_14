#!/bin/bash

eksctl delete cluster learnanddevops

rm ~/.kube/config
mv ~/.kube/work_config ~/.kube/config
