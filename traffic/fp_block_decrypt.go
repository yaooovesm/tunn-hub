package traffic

import "github.com/xtaci/kcp-go"

//
// BlockDecryptFP
// @Description:
//
type BlockDecryptFP struct {
	crypt kcp.BlockCrypt
}

//
// NewBlockDecryptFP
// @Description:
// @param key
// @return fp
// @return err
//
func NewBlockDecryptFP(crypt kcp.BlockCrypt) *BlockDecryptFP {
	return &BlockDecryptFP{
		crypt: crypt,
	}
}

//
// Init
// @Description:
// @receiver e
// @return bool
//
func (e *BlockDecryptFP) Init() bool {
	return true
}

//
// Process
// @Description:
// @receiver e
// @param raw
// @return []byte
//
func (e *BlockDecryptFP) Process(raw []byte) []byte {
	e.crypt.Decrypt(raw, raw)
	return raw
}
