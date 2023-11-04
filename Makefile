DOCKER_USERNAME := $(shell hack/get-docker-username.sh)

# =============================================================

build-server:
	docker build -t $(DOCKER_USERNAME)/node-latency-server:latest -f config/server/Dockerfile .
	docker push $(DOCKER_USERNAME)/node-latency-server:latest

run-server: build-server
	docker run -p 8080:8080 $(DOCKER_USERNAME)/node-latency-server:latest

deploy-server: build-server
	kubectl apply -f config/server/daemonset.yaml
	kubectl apply -f config/server/service.yaml

# =============================================================

build-client:
	docker build -t $(DOCKER_USERNAME)/node-latency-client:latest -f config/client/Dockerfile .
	docker push $(DOCKER_USERNAME)/node-latency-client:latest

run-client: build-client
	docker run -p 8081:8081 $(DOCKER_USERNAME)/node-latency-client:latest

deploy-client: build-client
	kubectl apply -f config/client/cluster-role.yaml
	kubectl apply -f config/client/cluster-role-binding.yaml
	kubectl apply -f config/client/service-account.yaml
	kubectl apply -f config/client/daemonset.yaml
	kubectl apply -f config/client/service.yaml

# =============================================================

test-env:
	./hack/create-cluster.sh