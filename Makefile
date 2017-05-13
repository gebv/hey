
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

tests:
	mkdir -p state/tarantool/_data
	mkdir -p state/tarantool/conf
	cp tarantool/app.lua state/tarantool/conf/
	cp tarantool/docker-compose.yml ./
	sudo docker-compose start tarantool
	TARANTOOL_SERVER=127.0.0.1:3301 TARANTOOL_USER_NAME=chrono TARANTOOL_USER_PASSWORD=chrono go test

clean:
	sudo docker-compose stop tarantool
	sudo rm -rf ./state
	rm docker-compose.yml
