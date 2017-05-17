all:
	mkdir -p state/tarantool/_data
	mkdir -p state/tarantool/conf
	cp ../../tarantool/app.lua state/tarantool/conf/
	cp ../../tarantool/docker-compose.yml ./
	-docker-compose up -d tarantool
	sleep 1
	TARANTOOL_SERVER=127.0.0.1:3301 TARANTOOL_USER_NAME=chrono TARANTOOL_USER_PASSWORD=chrono go run *.go

clean:
	sudo docker-compose stop tarantool
	sudo rm -rf ./state
	rm docker-compose.yml
