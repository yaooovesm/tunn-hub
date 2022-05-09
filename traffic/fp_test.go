package traffic

import (
	"bytes"
	"fmt"
	"github.com/gofrs/uuid"
	"github.com/xtaci/kcp-go"
	"testing"
	"time"
	"tunn-hub/logging"
)

type fp1 struct {
	test int
}

func (f fp1) Init() bool {
	return true
}

func (f fp1) Process(raw []byte) []byte {
	fmt.Println(f.test)
	return raw
}

func TestFlowProcessors(t *testing.T) {
	logging.Initialize()
	b := []byte{1}
	fp := NewFlowProcessor()
	fp.Register(fp1{test: 1}, "test1")
	fp.Register(fp1{test: 2}, "test2")
	fp.Register(fp1{test: 3}, "test3")
	fmt.Println("------------------------")
	fp.Process(b)
	fmt.Println("------------------------")
	//node := fp.GetByName("test3")
	//node.Processor.Process(b)
	fp.Register(fp1{test: 300}, "test3")
	fmt.Println("------------------------")
	fp.Process(b)
	fmt.Println("------------------------")
	fp.ProcessReverse(b)
	fmt.Println("------------------------")
	fp.Delete("test2")
	fp.Process(b)
	fmt.Println("------------------------")
	//node = fp.GetByName("test3")
	//node.Processor.Process([]byte{1})
	list := fp.List()
	fmt.Println(list)
}

func TestBlockEncryptFP(t *testing.T) {
	var key []byte
	p1, _ := uuid.NewV4()
	p2, _ := uuid.NewV4()
	key = append(key, p1.Bytes()...)
	key = append(key, p2.Bytes()...)
	c1, err := kcp.NewSalsa20BlockCrypt(key)
	if err != nil {
		fmt.Println(err)
		return
	}
	c2, err := kcp.NewSalsa20BlockCrypt(key)
	if err != nil {
		fmt.Println(err)
		return
	}

	fp1 := NewFlowProcessor()
	fp1.Register(NewBlockEncryptFP(c1), "encrypt")
	fp2 := NewFlowProcessor()
	fp2.Register(NewBlockDecryptFP(c2), "decrypt")
	time.Sleep(time.Second)

	for i := 0; i < 2; i++ {
		random, _ := uuid.NewV4()
		b := bytes.Repeat(random.Bytes(), 100)
		fmt.Println("origin  --> ", len(b))
		fp1.Process(b)
		fmt.Println("encrypt --> ", len(b))
		fp2.Process(b)
		fmt.Println("decrypt --> ", len(b))
	}
}
