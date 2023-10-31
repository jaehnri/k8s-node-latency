build-server:
	docker build -t jaehnri/node-server:latest -f config/server/Dockerfile .

run-server: build-server
	docker run -p 8080:8080 jaehnri/node-server:latest

deploy-server: build-server
	kubectl apply -f config/server/daemonset.yaml
	kubectl apply -f config/server/service.yaml

# =============================================================

build-client:
	docker build -t jaehnri/node-client:latest -f Dockerfile .

run-client: build-client
	docker run -p 8080:8080 jaehnri/node-client:latest

deploy-client: build-client
	kubectl apply -f config/client/daemonset.yaml
	kubectl apply -f config/client/service.yaml