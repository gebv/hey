package main

import (
	"api"
	"flag"
	// "github.com/golang/glog"
	_ "models"
	_ "store"
	"utils"

	"os"
	"os/signal"
	// _ "store"
	"syscall"
)

var flagConfigFile string 

func main() {
	flag.StringVar(&flagConfigFile, "config", "config.json", "")

	flag.Parse()

	utils.LoadConfig(flagConfigFile)

	api.NewServer()
	api.InitApi()
	api.StartServer()

	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	<-c

	api.StopServer()
}
