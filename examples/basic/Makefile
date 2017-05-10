all:
	sudo docker-compose start tarantool
	TARANTOOL_SERVER=127.0.0.1:3301 TARANTOOL_USER_NAME=chrono TARANTOOL_USER_PASSWORD=chrono go run *.go

clean:
	sudo docker-compose stop tarantool
	sudo rm ./state/tarantool/_data/*
