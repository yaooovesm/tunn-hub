package traffic

import (
	"github.com/klauspost/compress/zstd"
)

//
// ZSTDDecompressFP
// @Description:
//
type ZSTDDecompressFP struct {
	decoder *zstd.Decoder
}

//
// Init
// @Description:
// @receiver c
// @return bool
//
func (c *ZSTDDecompressFP) Init() bool {
	decoder, err := zstd.NewReader(nil, zstd.WithDecoderConcurrency(0))
	if err != nil {
		return false
	}
	c.decoder = decoder
	return true
}

//
// Process
// @Description:
// @receiver c
// @param raw
// @return []byte
//
func (c *ZSTDDecompressFP) Process(raw []byte) []byte {
	res, err := c.decoder.DecodeAll(raw, nil)
	if err != nil {
		return raw
	}
	return res
}

//
// Close
// @Description:
// @receiver c
//
func (c *ZSTDDecompressFP) Close() {
	c.decoder.Close()
}
