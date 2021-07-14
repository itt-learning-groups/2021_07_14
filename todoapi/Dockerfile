FROM golang:1.16-alpine AS builder
WORKDIR /go/src
COPY . .
RUN apk add --update make && \
    make build-alpine

FROM scratch
COPY --from=builder /go/src/cmd/todo_api_server/todo_api_server todo_api_server
USER 1001:1001
EXPOSE 8080
CMD ["./todo_api_server"]
