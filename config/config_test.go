package config

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestLoad(t *testing.T) {
	config := ServerConfigStorage{}
	config.ReadFromFile("E:\\TunnelTest\\server\\server.json")
	marshal, err := json.Marshal(config)
	if err != nil {
		return
	}
	fmt.Println(string(marshal))
}

func TestDump(t *testing.T) {
	Location = "E:\\TunnelTest\\server\\test.json"
	config := ServerConfigStorage{}
	config.ReadFromFile(Location)
	config.IPPool = IPPool{
		Start:   "192.168.20.2",
		End:     "192.168.20.100",
		Network: "192.168.20.0/24",
	}
	cfg := config.ToConfig()
	err := cfg.Storage()
	if err != nil {
		fmt.Println(err)
		return
	}
}
