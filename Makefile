build-server:
	docker build -t node-latency-server:test -f config/server/Dockerfile .

run-server: build-server
	docker run -p 8080:8080 node-latency-server:test

deploy-server: build-server
	kind load docker-image node-latency-server:test --name=kpng-e2e-ipv4-iptables
	kubectl apply -f config/server/deploy.yaml
	kubectl apply -f config/server/service.yaml

rollout-server: build-server
	kind load docker-image node-latency-server:test --name=kpng-e2e-ipv4-iptables
	kubectl rollout restart daemonset node-latency-server -n node-latency

# =============================================================

build-client:
	docker build -t node-latency-client:test -f config/client/Dockerfile .

run-client: build-client
	docker run -p 8081:8081 node-latency-client:test

deploy-client: build-client
	kind load docker-image node-latency-client:test --name=kpng-e2e-ipv4-iptables
	kubectl apply -f config/client/cluster-role.yaml
	kubectl apply -f config/client/cluster-role-binding.yaml
	kubectl apply -f config/client/service-account.yaml
	kubectl apply -f config/client/daemonset.yaml
	kubectl apply -f config/client/service.yaml

rollout-client: build-client
	kind load docker-image node-latency-client:test --name=kpng-e2e-ipv4-iptables
	kubectl rollout restart daemonset node-latency-client -n node-latency

# =============================================================

namespace:
	kubectl apply -f config/cluster/namespace.yaml

prometheus:
	helm install prometheus prometheus-community/prometheus -n node-latency -f config/prometheus/values.yaml

uninstall-prometheus:
	helm uninstall prometheus -n node-latency

grafana:
	helm install grafana grafana/grafana -n node-latency -f config/grafana/values.yaml

uninstall-grafana:
	helm uninstall grafana -n node-latency

node-latency:
	helm install node-latency config/node-latency/. -n node-latency

# =============================================================

test-env:
	./hack/create-cluster.sh

clean-test-env:
	kind delete cluster --name node-latency