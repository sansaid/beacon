GOBIN="$(shell go env GOPATH)/bin"
SWAGGER_DOWNLOAD_URL=$(shell curl -s https://api.github.com/repos/go-swagger/go-swagger/releases/latest | \
	jq -r '.assets[] | select(.name | contains("'"$(shell uname | tr '[:upper:]' '[:lower:]')"'_amd64")) | .browser_download_url')

${GOBIN}/golangci-lint: 
	@curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b ${GOBIN} v1.51.2

${GOBIN}/swagger:
	@curl -o ${GOBIN}/swagger -L'#' "${SWAGGER_DOWNLOAD_URL}"
	@chmod +x ${GOBIN}/swagger

${GOBIN}/swag:
	@go install github.com/swaggo/swag/cmd/swag@latest

${GOBIN}/mockgen:
	@go install github.com/golang/mock/mockgen@v1.6.0

# Used for generating swagger doc from Go annotations using `swag` and clients using `swagger`
.PHONY: swagger-gen
swagger-gen: ${GOBIN}/swag ${GOBIN}/swagger
	@swagger generate model -t beacond/ -f docs/swagger.yaml
	@swag fmt
	@swag init -g beacond/server/server.go
	@swagger generate client -t beacond/ -f docs/swagger.yaml
	@go mod tidy

.PHONY: deps
deps: ${GOBIN}/golangci-lint ${GOBIN}/swagger ${GOBIN}/mockgen
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