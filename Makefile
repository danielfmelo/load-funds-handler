img := danielfmelo/load-funds-handler:latest
wd := $(shell pwd)
cachevol=$(wd)/.gomodcachedir:/go/pkg/mod
rundocker := docker run --rm -v $(wd):/app -v $(cachevol) $(img)

image:
	docker build . -t $(img)

run: 
	go run cmd/load_funds_handler.go

docker-run: image docker-build
	$(rundocker) ./load_funds_handler 

build:
	go build -o ./load_funds_handler ./cmd/load_funds_handler.go

docker-build: image
	$(rundocker) go build -v -o ./load_funds_handler ./cmd/load_funds_handler.go

tests: 
	go test -timeout 20s -tags unit -race -coverprofile=coverage.out ./...

docker-tests: image
	$(rundocker) go test -timeout 20s -tags unit -race -coverprofile=coverage.out ./...

coverage: tests
	go tool cover -html=coverage.out -o=coverage.html
	xdg-open coverage.html
