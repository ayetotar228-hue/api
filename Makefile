.PHONY: build-api build-worker build-all

build-api:
	docker build \
		--build-arg SERVICE_NAME=api \
		--build-arg SERVICE_PATH=./cmd/api \
		--build-arg SERVICE_PORT=8080 \
		-t study-go/api:latest \
		-f Dockerfile .

build-worker:
	docker build \
		--build-arg SERVICE_NAME=worker \
		--build-arg SERVICE_PATH=./cmd/worker \
		--build-arg SERVICE_PORT=0 \
		-t study-go/worker:latest \
		-f Dockerfile .

build-all: build-api build-worker

deploy: build-all
	docker stack deploy -c docker-stack.yml study-go

clean:
	docker stack rm study-go
	docker rmi study-go/api:latest study-go/worker:latest