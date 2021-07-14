.PHONY: build
build:
	cd cmd/todo_api_server && make build

.PHONY: build-alpine
build-alpine:
	cd cmd/todo_api_server && make build-alpine

.PHONY: run
run:
	cd cmd/todo_api_server && make run

.PHONY: test test-unit
test test-unit:
	go test ./... | tee test-unit.log

.PHONY: test-coverage
test-coverage:
	go test -mod=vendor -coverprofile=cov.out ./...
	go tool cover -html=cov.out -o coverage.html

.PHONY: test-api-local
test-api-local:
	docker image build -t gotodoapi .
	docker container run --name gotodoapi-servicetest -p 8080:8080 --env SERVER_ADDR="" --env SERVER_PORT=8080 --rm -d gotodoapi
	newman run test/service/gotodoapi.postman_collection.json -e test/service/local.postman_environment.json
	docker container stop gotodoapi-servicetest

.PHONY: test-api-docker
test-api-docker:
	docker-compose -f docker-compose.test-service.yml up --build -d web
	docker-compose -f docker-compose.test-service.yml up test
	docker-compose -f docker-compose.test-service.yml down
