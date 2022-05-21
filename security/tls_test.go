package security

import (
	"fmt"
	"testing"
	"time"
	"tunn-hub/config"
)

func TestTLS(t *testing.T) {
	config.Location = "E:\\TunnelTest\\server\\test.json"
	config.Current.ReadFromFile(config.Location)
	certification := NewTunnX509Certification([]string{
		"127.0.0.1",
	}, []string{
		"localhost",
	}, time.Now().Add(time.Hour*86400))
	path := "./cert/"
	name, err := certification.CreateAndWriteWithRandomName("./cert/")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(name)
	config.Current.Security.CertPem = path + name + ".cert"
	config.Current.Security.KeyPem = path + name + ".key"
	err = config.Current.Storage()
	if err != nil {
		fmt.Println(err)
	}
}
