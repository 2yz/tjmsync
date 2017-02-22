package main

import (
	"io/ioutil"
	"log"
	"flag"
	"fmt"
	"os"
	"github.com/yezersky/tjmsync/lib"
)

var conf_path string

func main() {
	InitFlag()
	Init()
	worker := lib.NewWorker(lib.GetGlobalJobPool())
	worker.Start()
	server := &lib.StatusServer{}
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
	config, err2 := lib.ParseConfig(data)
	if err2 != nil {
		log.Fatal("Init err2: ", err2)
	}
	lib.InitGlobalConfig(config)
	lib.InitGlobalJobPool(config.Jobs)
	lib.InitJobLogger()
}