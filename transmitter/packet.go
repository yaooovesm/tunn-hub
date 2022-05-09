package transmitter

import (
	"errors"
)

var (
	ErrBadPacket = errors.New("bad packet")
)

//
// Packet
// @Description:
//
type Packet struct {
	Version  []byte
	Length   []byte
	Payload  []byte
	Checksum []byte
}

//
// Bytes
// @Description:
// @receiver p
// @return []byte
//
func (p *Packet) Bytes() []byte {
	b := make([]byte, 6)
	b[0] = p.Version[0]
	b[1] = p.Version[1]
	b[3] = p.Length[0]
	b[4] = p.Length[1]
	return append(b, p.Payload...)
}
