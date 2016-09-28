
vendor:
	go get -v github.com/stretchr/testify/assert \
		gopkg.in/jackc/pgx.v2 \
		github.com/satori/go.uuid
.PHONY: vendor

test: vendor
	GOPATH=${GOPATH}:${PWD} go test -v \
		-bench=. -benchmem \
		-run=. ./...
.PHONY: test
