kind create cluster --name kpng-test --config=config/cluster/kind.yaml

kubectl apply -f config/cluster/namespace.yaml

make deploy-server
make deploy-client
