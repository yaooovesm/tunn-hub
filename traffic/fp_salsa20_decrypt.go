package traffic

import "github.com/xtaci/kcp-go"

//
// Salsa20DecryptFP
// @Description:
//
type Salsa20DecryptFP struct {
	crypt kcp.BlockCrypt
}

//
// NewSalsa20DecryptFP
// @Description:
// @param key
// @return fp
// @return err
//
func NewSalsa20DecryptFP(key []byte) (fp *Salsa20DecryptFP, err error) {
	crypt, err := kcp.NewSalsa20BlockCrypt(key)
	if err != nil {
		return nil, err
	}
	return &Salsa20DecryptFP{
		crypt: crypt,
	}, nil
}

//
// Init
// @Description:
// @receiver e
// @return bool
//
func (e *Salsa20DecryptFP) Init() bool {
	return true
}

//
// Process
// @Description:
// @receiver e
// @param raw
// @return []byte
//
func (e *Salsa20DecryptFP) Process(raw []byte) []byte {
	//数据长度要大于8
	if len(raw) <= 8 {
		return raw
	}
	e.crypt.Decrypt(raw, raw)
	return raw
}

//
// Close
// @Description:
// @receiver e
//
func (e *Salsa20DecryptFP) Close() {
}
