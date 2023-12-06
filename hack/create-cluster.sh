kind create cluster --name node-latency --config=config/cluster/kind.yaml

kubectl apply -f config/cluster/namespace.yaml

helm install node-latency config/node-latency/. -n node-latency

make prometheus
make grafana
