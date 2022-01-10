# Pinot Playground
Simple project to play with Apache Pinot:
 - installation in k8s
 - installation of cert-manager
 - monitor certificates expiration

All the deployments are done in minikube on MacOS.

## Install

### Minikube
Follow instruction to install minikube provided at https://minikube.sigs.k8s.io/docs/start/
To be able to interact with Docker running in minikube VM:
```shell
eval $(minikube docker-env)
```

### Helm
If helm is not installed on your system
```shell
brew install helm
```

### Apache Pinot
Following Apache Pinot documentation:
```shell
helm repo add pinot https://raw.githubusercontent.com/apache/pinot/master/kubernetes/helm
kubectl create ns pinot-quickstart
helm install pinot pinot/pinot -n pinot-quickstart --set cluster.name=pinot --set server.replicaCount=2
```
To verify the deployment
```shell
kubectl -n pinot-quickstart get all
```

For more details, go to https://docs.pinot.apache.org/basics/getting-started/kubernetes-quickstart#2.-setting-up-a-pinot-cluster-in-kubernetes

### Cert-manager
Install cert-manager with regular manifest
```shell
kubectl apply -f https://github.com/jetstack/cert-manager/releases/download/v1.6.1/cert-manager.yaml
```

To verify the deployment
```shell
kubectl -n cert-manager get all
```

For more details, go to https://cert-manager.io/docs/installation/

### Cert-monitor
This section is related to the cert-monitor that will list all certificate CRD to verify their expiration.
The cert-monitor is deployed as a k8s cron job that runs every minute.
All commands described in that section must be run in the cert-monitor directory.
```shell
cd cert-monitor
```

#### Build

```shell
make docker
```

#### Deployment
```shell
kubectl create ns cert-monitor
kubectl -n cert-monitor apply -f kubernetes/rbac.yml
kubectl -n cert-monitor apply -f kubernetes/config.yml
kubectl -n cert-monitor apply -f kubernetes/job.yml
```
To see the logs of last job
```shell
kubectl -n cert-monitor logs `kubectl -n cert-monitor get po | grep cert-monitor | tail -1 | awk '{print $1}'`
```

#### Creating certificate