GOBIN=$(shell go env GOPATH)


${GOBIN}/bin/golangci-lint: 
	@curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b ${GOBIN}/bin v1.51.2

.PHONY: swagger-gen
swagger-gen:
	@go install github.com/swaggo/swag/cmd/swag@latest
	@swag fmt
	@swag init -g beacond/server/server.go

.PHONY: deps
deps: ${GOBIN}/bin/golangci-lint
	@go install github.com/golang/mock/mockgen@v1.6.0
	@go get -d ./...

.PHONY: mocks
mocks: deps
	@go generate ./...

.PHONY: test
test: deps mocks
	@golangci-lint run ./...
	@go test ./...

.PHONY: build
build:
	@go build .