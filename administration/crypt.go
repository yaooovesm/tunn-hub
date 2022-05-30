package administration

import (
	"bytes"
	"encoding/base64"
	"github.com/xtaci/kcp-go"
	"tunn-hub/traffic"
)

type Crypt struct {
	encrypt *traffic.FlowProcessors
	decrypt *traffic.FlowProcessors
}

//
// NewCrypt
// @Description:
// @param key
// @return crypt
// @return err
//
func NewCrypt(key []byte) (crypt Crypt, err error) {
	if len(key) > 16 {
		key = key[:16]
	} else if len(key) < 16 {
		padding := bytes.Repeat([]byte{0}, 16-len(key))
		key = append(key, padding...)
	}
	aes, err := kcp.NewAESBlockCrypt(key)
	if err != nil {
		return Crypt{}, err
	}

	encrypt := traffic.NewFlowProcessor()
	encrypt.Name = "orm_data_encrypt_processor"
	encrypt.Register(traffic.NewBlockEncryptFP(aes), "encrypt")
	encrypt.Register(&traffic.ZSTDCompressFP{}, "compress")

	decrypt := traffic.NewFlowProcessor()
	decrypt.Name = "orm_data_decrypt_processor"
	decrypt.Register(&traffic.ZSTDDecompressFP{}, "decompress")
	decrypt.Register(traffic.NewBlockDecryptFP(aes), "decrypt")

	return Crypt{
		encrypt: encrypt,
		decrypt: decrypt,
	}, nil
}

//
// Encrypt
// @Description:
// @receiver c
// @param data
// @return string
//
func (c Crypt) Encrypt(data string) string {
	return base64.StdEncoding.EncodeToString(c.encrypt.Process([]byte(data)))
}

//
// Decrypt
// @Description:
// @receiver c
// @param data
// @return string
//
func (c Crypt) Decrypt(data string) string {
	raw, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		return ""
	}
	return string(c.decrypt.Process(raw))
}
