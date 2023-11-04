kind create cluster --name kpng-test --config=config/cluster/kind.yaml

make deploy-server
make deploy-client
