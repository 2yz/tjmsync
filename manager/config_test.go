package manager

import (
	"fmt"
	"io/ioutil"
	"testing"
)

func TestParseConfig(t *testing.T) {
	data, err1 := ioutil.ReadFile("../example/tjmsync.toml")
	if err1 != nil {
		t.Error(err1)
		return
	}
	config, err2 := ParseConfig(data)
	if err2 != nil {
		t.Error("Config Parse Error: ", err2)
		return
	}
	fmt.Println("Config Jobs: ", config.Jobs)
	fmt.Println("Config Log: ", config.Log)
	fmt.Println("Config Status Server: ", config.StatusServer)
}
