download:
	go mod download

build:
	go build -v ./...

test:
	go test -race -v ./...

lint:
	golangci-lint run --timeout 60s --max-same-issues 50 ./...

lint-fix:
	golangci-lint run --timeout 60s --max-same-issues 50 --fix ./...
