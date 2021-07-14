.PHONY: build
build:
	make clean
	go build -mod=vendor -o todo_api_server main.go

.PHONY: build-alpine
build-alpine:
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -mod=vendor -a -ldflags="-w -extldflags '-static'" -o todo_api_server main.go

.PHONY: run
run:
	go run main.go

.PHONY: clean
clean:
	rm -f main
	rm -f todo_api_server
