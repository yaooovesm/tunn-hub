package transmitter

import (
	"bytes"
	"encoding/binary"
	"hash/crc32"
	"io"
	"net"
)

type TunReader struct {
	conn    io.ReadWriter
	version Version
	header  []byte
	convbuf *bytes.Buffer
	crcbuf  *bytes.Buffer
}

//
// NewTunReader
// @Description:
// @param conn
// @return *TunReader
//
func NewTunReader(conn io.ReadWriter, ver Version) *TunReader {
	return &TunReader{
		conn:    conn,
		version: ver,
		header:  make([]byte, 6),
		convbuf: bytes.NewBuffer([]byte{}),
		crcbuf:  bytes.NewBuffer([]byte{}),
	}
}

//
// readSizeFromHeader
// @Description:
// @receiver r
// @return size
//
func (r *TunReader) readSizeFromHeader() (size uint16) {
	r.convbuf.Write([]byte{r.header[3], r.header[4]})
	_ = binary.Read(r.convbuf, binary.BigEndian, &size)
	return
}

//
// crc32
// @Description:
// @receiver r
// @param b
// @return []byte
//
func (r *TunReader) crc32(b []byte) uint32 {
	return crc32.ChecksumIEEE(b)
}

//
// Read
// @Description:
// @receiver r
// @param p
// @return n
// @return err
//
func (r *TunReader) Read() (p []byte, err error) {
	n, err := r.conn.Read(r.header)
	if err != nil {
		return
	}
	if n != 6 {
		return nil, ErrBadPacket
	}
	if r.header[0] != r.version[0] || r.header[1] != r.version[1] {
		return nil, ErrBadPacket
	}
	ps := r.readSizeFromHeader()
	switch r.version {
	case V1:
		p = make([]byte, ps)
		n, err = r.conn.Read(p)
		if err != nil {
			return nil, err
		}
	case V2:
		pscrc := ps + 4
		buf := make([]byte, pscrc)
		n, err = r.conn.Read(buf)
		if err != nil {
			return nil, err
		}
		if uint16(n) != pscrc {
			return nil, ErrBadPacket
		}
		defer func() {
			r.convbuf.Reset()
		}()
		//checksum : crc32
		p = buf[:ps]
		var checksum uint32
		r.convbuf.Write(buf[ps:])
		_ = binary.Read(r.convbuf, binary.BigEndian, &checksum)
		if crc32.ChecksumIEEE(append(r.header, p...)) != checksum {
			p = []byte{}
			return nil, ErrBadPacket
		}
	}
	return
}

//
// Conn
// @Description:
// @receiver r
// @return io.ReadWriter
//
func (r *TunReader) Conn() net.Conn {
	return r.conn.(net.Conn)
}
