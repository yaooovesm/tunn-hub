package authentication

const (
	MaxAuthPacketSize = 4096
)

type PacketType int

const (
	PacketTypeUnknown = 255
	PacketTypeMsg     = 0
	PacketTypeReport  = 1
	PacketTypeLogin   = 10
	PacketTypeLogout  = 11
	PacketTypeKick    = 12
	PacketTypeRestart = 13
)

//
// TransportPacket
// @Description:
//
type TransportPacket struct {
	Type    PacketType
	UUID    string
	Payload []byte
}

//
// NewTransportPacket
// @Description:
// @return *TransportPacket
//
func NewTransportPacket() *TransportPacket {
	return &TransportPacket{
		Type: PacketTypeUnknown,
	}
}

//
// Encode
// @Description:
// @receiver tp
// @return []byte
//
func (tp *TransportPacket) Encode() []byte {
	var b []byte
	b = append(b, byte(tp.Type))
	b = append(b, []byte(tp.UUID)...)
	b = append(b, tp.Payload...)
	return b
}

//
// Decode
// @Description:
// @receiver tp
// @return []byte
//
func (tp *TransportPacket) Decode(bytes []byte) (err error) {
	tp.Type = PacketType(bytes[0])
	tp.UUID = string(bytes[1:33])
	tp.Payload = bytes[33:]
	return nil
}
