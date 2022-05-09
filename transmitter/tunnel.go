package transmitter

import (
	"net"
	"time"
)

type Tunnel struct {
	reader *TunReader
	writer *TunWriter
	ver    Version
	conn   net.Conn
}

//
// NewTunnel
// @Description:
// @param conn
// @param version
// @return *Tunnel
//
func NewTunnel(conn net.Conn, version Version) *Tunnel {
	return &Tunnel{
		reader: NewTunReader(conn, version),
		writer: NewTunWriter(conn, version),
		ver:    version,
		conn:   conn,
	}
}

//
// Read
// @Description:
// @receiver t
// @return pl
// @return err
//
func (t *Tunnel) Read() (pl []byte, err error) {
	return t.reader.Read()
}

//
// Write
// @Description:
// @receiver t
// @param pl
// @return int
// @return error
//
func (t *Tunnel) Write(pl []byte) (int, error) {
	return t.writer.Write(pl)
}

//
// Close
// @Description:
// @receiver t
// @return error
//
func (t *Tunnel) Close() error {
	return t.conn.Close()
}

//
// SetDeadline
// @Description:
// @receiver t
// @param time
// @return error
//
func (t *Tunnel) SetDeadline(time time.Time) error {
	return t.conn.SetDeadline(time)
}

//
// RemoteAddr
// @Description:
// @receiver t
// @return net.Addr
//
func (t *Tunnel) RemoteAddr() net.Addr {
	return t.conn.RemoteAddr()
}
