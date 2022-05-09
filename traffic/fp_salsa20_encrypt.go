package traffic

import "github.com/xtaci/kcp-go"

//
// Salsa20EncryptFP
// @Description:
//
type Salsa20EncryptFP struct {
	crypt kcp.BlockCrypt
}

//
// NewSalsa20EncryptFP
// @Description:
// @param key
// @return fp
// @return err
//
func NewSalsa20EncryptFP(key []byte) (fp *Salsa20EncryptFP, err error) {
	crypt, err := kcp.NewSalsa20BlockCrypt(key)
	if err != nil {
		return nil, err
	}
	return &Salsa20EncryptFP{
		crypt: crypt,
	}, nil
}

//
// Init
// @Description:
// @receiver e
// @return bool
//
func (e *Salsa20EncryptFP) Init() bool {
	return true
}

//
// Process
// @Description:
// @receiver e
// @param raw
// @return []byte
//
func (e *Salsa20EncryptFP) Process(raw []byte) []byte {
	//数据长度要大于8
	if len(raw) <= 8 {
		return raw
	}
	e.crypt.Encrypt(raw, raw)
	return raw
}
