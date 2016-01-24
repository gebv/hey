package main

import (
	_ "api"
	"flag"
	"github.com/golang/glog"
	_ "models"
	_ "store"
	_ "utils"
)

func main() {
	flag.Parse()
	glog.Infof("Start...")
}
