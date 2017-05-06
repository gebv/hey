package tarantool

import (
	"os"
	"time"

	tarantool "github.com/tarantool/go-tarantool"
)

var (
	DefaultTimeout   time.Duration = 500 * time.Millisecond
	DefaultReconnect time.Duration = 1 * time.Second
)

func SetupFromENV() (client *tarantool.Connection, err error) {
	// 127.0.0.1:3013
	server := os.Getenv("TARANTOOL_SERVER")
	user := os.Getenv("TARANTOOL_USER_NAME")
	pwd := os.Getenv("TARANTOOL_USER_PASSWORD")

	opts := tarantool.Opts{
		Timeout:       DefaultTimeout,
		Reconnect:     DefaultReconnect,
		MaxReconnects: 3,
		User:          user,
		Pass:          pwd,
	}

	client, err = tarantool.Connect(server, opts)
	if err != nil {
		return client, err
	}

	_, err = client.Ping()

	return client, err
}
