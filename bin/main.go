package main

import (
	"io/ioutil"
	"log"
	"github.com/yezersky/tjmsync"
	"flag"
	"fmt"
	"os"
)

var conf_path string

func main() {
	InitFlag()
	Init()
	worker := tjmsync.NewWorker(tjmsync.GetGlobalJobPool())
	worker.Start()
	server := &tjmsync.StatusServer{}
	server.Serve()
}

func InitFlag() {
	flag.StringVar(&conf_path, "conf", "" , "config file path")
	flag.Parse()
	if conf_path == "" {
		fmt.Println("require -conf param")
		os.Exit(1)
	}
}

func Init() {
	data, err1 := ioutil.ReadFile(conf_path)
	if err1 != nil {
		log.Fatal("Init err1: ", err1)
	}
	config, err2 := tjmsync.ParseConfig(data)
	if err2 != nil {
		log.Fatal("Init err2: ", err2)
	}
	tjmsync.InitGlobalConfig(config)
	tjmsync.InitGlobalJobPool(config.Jobs)
	tjmsync.InitJobLogger()
}