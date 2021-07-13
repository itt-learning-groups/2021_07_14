# Kubernetes "Lab" Exercise: Deploy a web api -- *continued from last week (and somewhat modified/extended)*

## Goals

* Pick up from where we left off last week (roughly [here](https://github.com/us-learn-and-devops/2021_07_07#5-scale-the-deployment), at part "5" of last week's lab).
* Diagnose and solve our scaling limitation.
* Understand resource requests and limitations.
* Understand liveness/readiness probes.
* Understand deployment strategies (e.g. "rollout") and deployment rollback.
* Understand how to troubleshoot a couple of failure scenarios.
* Understand the k8s object/concept landscape at a high level, our roadmap through these for the next 10-12 sessions, and where we're at on that roadmap so far.

## Lab

### 1. Setup

* Clone the updated `todoapi` project here.
* Add AWS creds to your terminal session and start up your cluster. You can use the `cluster-deploy.sh` script if you like.
* Deploy the todoapi `deployment` and `service` from where we left off last week. You can use the `k8s/deploy.sh` script if you like.

### 2. Scaling Limitations

* Try scaling the deployment's RelicaSet to 4 replicas:
  * Best to use the `spec.replicas` in your `deployment.yaml` and re-apply the config
  * But you can also use kubectl directly: [See the docs](https://kubernetes.io/docs/concepts/workloads/controllers/deployment/#scaling-a-deployment)
* [Watch or list the pods](https://kubernetes.io/docs/reference/kubectl/cheatsheet/#viewing-finding-resources) from your deployment and make sure you see your single pod become 4 pods. This should happen very quickly: The todoapi is simple.
* Now try scaling to 5 replicas and again watch or list the pods. This time you should see the fifth pod get stuck in "pending" state.

What's happening?

The most typical reason to see a pod get stuck in "pending" like this is [insufficient resources](https://kubernetes.io/docs/tasks/debug-application-cluster/debug-application/#my-pod-stays-pending).

