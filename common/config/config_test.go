package config

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestLoad(t *testing.T) {
	Current.ReadFromFile("D:\\code\\golang\\tunnel\\common\\config\\tunnel_client.json")
	marshal, err := json.Marshal(Current)
	if err != nil {
		return
	}
	fmt.Println(string(marshal))
}
