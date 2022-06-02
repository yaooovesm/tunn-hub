package traffic

import "github.com/xtaci/kcp-go"

//
// BlockEncryptFP
// @Description:
//
type BlockEncryptFP struct {
	crypt kcp.BlockCrypt
}

//
// NewBlockEncryptFP
// @Description:
// @param key
// @return fp
// @return err
//
func NewBlockEncryptFP(crypt kcp.BlockCrypt) *BlockEncryptFP {
	return &BlockEncryptFP{
		crypt: crypt,
	}
}

//
// Init
// @Description:
// @receiver e
// @return bool
//
func (e *BlockEncryptFP) Init() bool {
	return true
}

//
// Process
// @Description:
// @receiver e
// @param raw
// @return []byte
//
func (e *BlockEncryptFP) Process(raw []byte) []byte {
	e.crypt.Encrypt(raw, raw)
	return raw
}

//
// Close
// @Description:
// @receiver e
//
func (e *BlockEncryptFP) Close() {
}
