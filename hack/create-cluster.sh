kind create cluster --name node-latency --config=config/cluster/kind.yaml

kubectl apply -f config/cluster/namespace.yaml

make deploy-server
make deploy-client

make prometheus
make grafana
