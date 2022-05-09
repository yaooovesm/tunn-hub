package administration

import (
	"fmt"
	"testing"
	"tunn-hub/common/logging"
)

func TestCrypt(t *testing.T) {
	logging.Initialize()
	crypt, err := NewCrypt([]byte("0123456789abcdef"))
	if err != nil {
		fmt.Println(err)
		return
	}
	encrypt := crypt.Encrypt("123")
	fmt.Println("encrypt --> ", encrypt)
	decrypt := crypt.Decrypt(encrypt)
	fmt.Println("decrypt --> ", decrypt)
}
