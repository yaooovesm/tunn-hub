package traffic

import (
	"github.com/klauspost/compress/zstd"
)

//
// ZSTDCompressFP
// @Description:
//
type ZSTDCompressFP struct {
	encoder *zstd.Encoder
}

//
// Init
// @Description:
// @receiver c
// @return bool
//
func (c *ZSTDCompressFP) Init() bool {
	encoder, err := zstd.NewWriter(nil, zstd.WithEncoderLevel(zstd.SpeedBetterCompression))
	if err != nil {
		return false
	}
	c.encoder = encoder
	return true
}

//
// Process
// @Description:
// @receiver c
// @param raw
// @return []byte
//
func (c *ZSTDCompressFP) Process(raw []byte) []byte {
	res := c.encoder.EncodeAll(raw, make([]byte, 0, len(raw)))
	return res
}

//
// Close
// @Description:
// @receiver c
//
func (c *ZSTDCompressFP) Close() {
	_ = c.encoder.Close()
}
