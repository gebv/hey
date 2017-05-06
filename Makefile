
vendor:
	go get -v github.com/stretchr/testify/assert
	go get -v github.com/satori/go.uuid
	go get -v github.com/tarantool/go-tarantool
.PHONY: vendor

test: vendor
	go test -v \
		-bench=. -benchmem \
		-run=. ./examples/...
.PHONY: test
