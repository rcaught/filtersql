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

release:
	sh -c 'git tag ${TAG} && git push origin $TAG && GOPROXY=proxy.golang.org go list -m github.com/rcaught/filtersql@${TAG}'