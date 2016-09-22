
BUILD_DATE = `date -u +%Y-%m-%dT%H:%M:%S%z`
BUILD_HASH = `git rev-parse HEAD 2>/dev/null || echo "???"`
BUILD_ID = ${CI_BUILD_ID}  
NAME = ${CI_PROJECT_NAME}

FLAGS ?= -a --installsuffix cgo -ldflags \
    "-s -X 'config.BuildDate=$(BUILD_DATE)' \
    -X config.BuildHash=$(BUILD_HASH) \
    -X config.BuildID=$(BUILD_ID)" \
    src/cmd/p356/main.go

build: vendor
	GOOS="linux" GOARCH="amd64" GOPATH=${GOPATH}:${PWD} CGO_ENABLED=0 go build -o bin/${NAME} -v ${FLAGS}
.PHONY: build

vendor:
	go get -v github.com/stretchr/testify/assert \
		gopkg.in/jackc/pgx.v2 \
		github.com/satori/go.uuid
.PHONY: vendor

test: vendor
	GOPATH=${GOPATH}:${PWD} go test -v \
		-bench=. -benchmem \
		-run=. ./src/...
.PHONY: test

dev-build:
	docker build -t "local/${NAME}" -f dev.Dockerfile .
.PHONY: dev-build