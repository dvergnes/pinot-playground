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

### Kafka
Kafka is used to collect the alerts emitted by the cert-monitor. The topic is cert-monitor-alerts.
Install Kafka
```shell
helm repo add incubator https://charts.helm.sh/incubator
helm install -n pinot-quickstart kafka incubator/kafka --set replicas=1,zookeeper.image.tag=latest
```

To verify the deployment
```shell
kubectl get all -n pinot-quickstart | grep kafka
```

To create the cert-monitor-alerts topic
```shell
kubectl -n pinot-quickstart exec kafka-0 -- kafka-topics --zookeeper kafka-zookeeper:2181 --topic cert-monitor-alerts --create --partitions 1 --replication-factor 1
```

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

To create Issuers
```shell
kubectl apply -f cert-manager/kubernetes/issuers.yml
```

To create a certificate. The created certificate will expire every hour. It is created in the sandbox namespace with the
name vergnes-com.
```shell
kubectl apply -f cert-manager/kubernetes/certificate.yml
```

#### Simulate cert-manager outage
The cert-manager deployment is responsible of renewing the certificates before their expiration. To verify that the
cert-monitor can detect expired certificates we can simulate that the cert-manager is down by setting the number of
replicas to 0.

To simulate that the cert-manager is down
```shell
kubectl -n cert-manager scale deploy/cert-manager --replicas=0
```

To restore the deployment of the cert-manager
```shell
kubectl -n cert-manager scale deploy/cert-manager --replicas=1
```

### Cert-monitor
This section is related to the cert-monitor that will list all certificate CRD to verify their expiration.
The cert-monitor is deployed as a k8s cron job that runs every minute.
The cert-monitor generates alerts as messages in the kafka topic cert-monitor-alerts. The alert is JSON encoded, it
contains: level, message, certificate location (name, namespace), timestamp, pod name who generated the alert. 

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

