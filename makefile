docker:
	@echo Building Docker Image
	@docker build -t nitricimages/dev-api-gateway .

clean:
	@echo Cleaning build artifacts
	@rm -rf bin/

install:
	@echo Fetching go dependencies
	@go mod download

build:
	@echo Building Go app
	@CGO_ENABLED=0 GOOS=linux go build -o bin/api-gateway -ldflags="-extldflags=-static" main.go
