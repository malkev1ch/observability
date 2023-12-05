# **Deployment**

Contains the order of the valid deployment. The command must be executed from the root of the project.

## Installing the Charts

Jaeger is deployed with k8s operator via helm chart
[Link1](https://github.com/jaegertracing/helm-charts/tree/main/charts/jaeger-operator)
[Link2](https://www.jaegertracing.io/docs/1.50/operator).

Opentelemetry collector is deployed without k8s operator via helm chart. [Link1](https://github.com/open-telemetry/opentelemetry-helm-charts/tree/main/charts/opentelemetry-collector).

**Step: 1** First, install the cert-manager on the k8s cluster.
```console
$ kubectl apply -f https://github.com/cert-manager/cert-manager/releases/download/v1.13.2/cert-manager.yaml
```

**Step: 2** Add the Jaeger Tracing Helm repository:
```console
$ helm repo add jaegertracing https://jaegertracing.github.io/helm-charts
```
You can then run helm search repo jaegertracing to see the charts.

**Step: 3** To install the jaeger chart with the release name `jaeger` in `jaeger` namespace:
```console
$ kubectl create namespace jaeger
$ helm install jaeger jaegertracing/jaeger-operator -n jaeger
```

**Step: 4** To install the default jaeger AllInOne strategy:
```console
$ kubectl apply -f deploy/jaeger/simplest.yaml
```

**Step: 5** To install the opentelemetry chart with the release name `otel-collector` in `otel-collector` namespace:
```console
$ kubectl create namespace otel-collector
$ helm install otel-collector-ds open-telemetry/opentelemetry-collector --values deploy/otel/daemonset.yaml -n otel-collector
```

## Installing Ingress
Before installing ingress verify the components running fine for jaeger and opentelemetry or not.
```console
$ kubectl get po -n jaeger
$ kubectl get svc -n jaeger
$ kubectl get po -n otel-collector
$ kubectl get svc -n otel-collector
```


## Uninstalling the Charts

To uninstall/delete the charts:
```console
$ helm delete jaeger
$ helm delete otel-collector
```
