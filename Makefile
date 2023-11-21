build-server:
	docker build -t node-latency-server:test -f config/server/Dockerfile .

run-server: build-server
	docker run -p 8080:8080 node-latency-server:test

deploy-server: build-server
	kind load docker-image node-latency-server:test --name=node-latency
	kubectl apply -f config/server/daemonset.yaml
	kubectl apply -f config/server/service.yaml

rollout-server: build-server
	kind load docker-image node-latency-server:test --name=node-latency
	kubectl rollout restart daemonset node-latency-server -n node-latency

# =============================================================

build-client:
	docker build -t node-latency-client:test -f config/client/Dockerfile .

run-client: build-client
	docker run -p 8081:8081 node-latency-client:test

deploy-client: build-client
	kind load docker-image node-latency-client:test --name=node-latency
	kubectl apply -f config/client/cluster-role.yaml
	kubectl apply -f config/client/cluster-role-binding.yaml
	kubectl apply -f config/client/service-account.yaml
	kubectl apply -f config/client/daemonset.yaml
	kubectl apply -f config/client/service.yaml

rollout-client: build-client
	kind load docker-image node-latency-client:test --name=node-latency
	kubectl rollout restart daemonset node-latency-client -n node-latency

# =============================================================

prometheus:
	helm install prometheus prometheus-community/prometheus -n node-latency -f config/prometheus/values.yaml

uninstall-prometheus:
	helm uninstall prometheus -n node-latency

test-env:
	./hack/create-cluster.sh