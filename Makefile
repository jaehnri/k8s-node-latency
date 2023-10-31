build-server:
	docker build -t jaehnri/node-server:latest -f Dockerfile .

run-server: build-server
	docker run -p 8080:8080 jaehnri/node-server:latest

deploy-server: build-server
	kubectl apply -f config/server/daemonset.yaml
	kubectl apply -f config/server/service.yaml