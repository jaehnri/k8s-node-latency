kind create cluster --name node-latency --config=config/cluster/kind.yaml

kubectl apply -f config/cluster/namespace.yaml

make node-latency
make prometheus
make grafana
