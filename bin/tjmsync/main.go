package main

import (
	"flag"
	"log"
	"io/ioutil"
	"github.com/yezersky/tjmsync/manager"
	"github.com/yezersky/grpool"
)

var (
	conf_path string = "/etc/tjmsync.toml"
)

func main() {
	flag.StringVar(&conf_path, "conf", "", "config file path. Default: /etc/tjmsync.toml")
	flag.Parse()

	config, err := LoadConfig()
	if err != nil {
		log.Fatal("LoadConfig Fatal Error:", err)
	}
	manager.SetGlobalConfig(config)
	manager.SetJobLogger(config.Log)
	wp, err := grpool.NewWorkerPool(config.MaxWorkerNumber)
	if err != nil {
		log.Fatal("NewWorkerPool Fatal Error:", err)
	}
	wp.Start()
	manager.SetGlobalWorkerPool(wp)
	m := manager.NewManager(config.Jobs, wp)
	err = m.Start()
	if err != nil {
		log.Fatal("Manager Start Fatal Error:", err)
	}
	manager.SetGlobalManager(m)

	server := &manager.StatusServer{}
	server.Serve()
}

func LoadConfig() (*manager.Config, error) {
	data, err := ioutil.ReadFile(conf_path)
	if err != nil {
		return nil, err
	}
	return manager.ParseConfig(data)
}
