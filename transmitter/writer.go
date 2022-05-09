package transmitter

import (
	"bytes"
	"encoding/binary"
	"hash/crc32"
	"io"
	"net"
)

//
// TunWriter
// @Description:
//
type TunWriter struct {
	conn    io.ReadWriter
	version Version
	vb      []byte
	header  []byte
	convbuf *bytes.Buffer
	crcbuf  *bytes.Buffer
}

//
// NewTunWriter
// @Description:
// @param conn
// @param ver
// @return *TunWriter
//
func NewTunWriter(conn io.ReadWriter, ver Version) *TunWriter {
	return &TunWriter{
		conn:    conn,
		version: ver,
		vb:      []byte(ver),
		header:  make([]byte, 6),
		convbuf: bytes.NewBuffer([]byte{}),
		crcbuf:  bytes.NewBuffer([]byte{}),
	}
}

//
// crc32
// @Description:
// @receiver w
// @param b
// @return []byte
//
func (w TunWriter) crc32(b []byte) []byte {
	defer func() {
		w.crcbuf.Reset()
	}()
	val := crc32.ChecksumIEEE(b)
	err := binary.Write(w.crcbuf, binary.BigEndian, &val)
	if err != nil {
		return nil
	}
	return w.crcbuf.Bytes()
}

//
// Write
// @Description:
// @receiver w
// @param p
// @return n
// @return err
//
func (w *TunWriter) Write(p []byte) (n int, err error) {
	if l := len(p); l > 0 {
		length := uint16(l)
		_ = binary.Write(w.convbuf, binary.BigEndian, &length)
		pkt := &Packet{
			Version: w.vb,
			Length:  w.convbuf.Bytes(),
			Payload: p,
		}
		w.convbuf.Reset()
		switch w.version {
		case V1:
			n, err = w.conn.Write(pkt.Bytes())
		case V2:
			b := pkt.Bytes()
			n, err = w.conn.Write(append(b, w.crc32(b)...))
		default:
			n, err = w.conn.Write(pkt.Bytes())
		}
	}
	return
}

//
// Conn
// @Description:
// @receiver w
// @return net.Conn
//
func (w *TunWriter) Conn() net.Conn {
	return w.conn.(net.Conn)
}
