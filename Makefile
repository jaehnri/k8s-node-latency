DOCKER_USERNAME := $(shell docker info | sed '/Username:/!d;s/.* //')

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
	docker build -t $(DOCKER_USERNAME)/node-latency-client:latest -f Dockerfile .

run-client: build-client
	docker run -p 8080:8080 $(DOCKER_USERNAME)/node-latency-client:latest
