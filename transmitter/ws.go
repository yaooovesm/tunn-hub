package transmitter

import (
	"bytes"
	"errors"
	"github.com/gorilla/websocket"
	"io"
	"net"
	"time"
)

//
// WrapWSConn
// @Description:
// @param wsconn
// @return *WSConn
//
func WrapWSConn(wsconn *websocket.Conn) *WSConn {
	return &WSConn{
		wsconn: wsconn,
		buf:    bytes.NewBuffer([]byte{}),
		closed: false,
	}
}

//
// WSConn
// @Description:
//
type WSConn struct {
	wsconn *websocket.Conn
	buf    *bytes.Buffer
	closed bool
}

//
// readToBuf
// @Description:
// @receiver w
// @return n
// @return err
//
func (w *WSConn) readToBuf() (n int, err error) {
	_, p, err := w.wsconn.ReadMessage()
	if err != nil {
		return 0, err
	}
	write, err := w.buf.Write(p)
	return write, err
}

var errClosed = errors.New("connection closed")

//
// Read
// @Description:
// @receiver w
// @param b
// @return n
// @return err
//
func (w *WSConn) Read(b []byte) (n int, err error) {
	length := len(b)
	if w.closed {
		return 0, errClosed
	}
	if w.buf.Len() == 0 {
		defer func() {
			if err := recover(); err != nil {
				w.closed = true
			}
		}()
		_, err = w.readToBuf()
		if err != nil {
			return 0, err
		}
	}
	if w.buf.Len() == 0 {
		_, _ = w.readToBuf()
	}
	if length <= 0 {
		return
	} else {
		//tmp := make([]byte, length)
		read, err := w.buf.Read(b)
		if err == io.EOF || read == 0 {
			return 0, nil
		}
		if err != nil {
			_ = w.Close()
			return 0, err
		}
		//for i := 0; i < read; i++ {
		//	b[i] = tmp[i]
		//}
		//copy(b, tmp)
		return read, nil
	}
}

//
// Write
// @Description:
// @receiver w
// @param b
// @return n
// @return err
//
func (w *WSConn) Write(b []byte) (n int, err error) {
	err = w.wsconn.WriteMessage(websocket.BinaryMessage, b)
	if err != nil {
		return 0, err
	}
	return len(b), nil
}

//
// Close
// @Description:
// @receiver w
// @return error
//
func (w *WSConn) Close() error {
	w.closed = true
	return w.wsconn.Close()
}

//
// LocalAddr
// @Description:
// @receiver w
// @return net.Addr
//
func (w WSConn) LocalAddr() net.Addr {
	return w.wsconn.LocalAddr()
}

//
// RemoteAddr
// @Description:
// @receiver w
// @return net.Addr
//
func (w WSConn) RemoteAddr() net.Addr {
	return w.wsconn.RemoteAddr()
}

//
// SetDeadline
// @Description:
// @receiver w
// @param t
// @return err
//
func (w *WSConn) SetDeadline(t time.Time) (err error) {
	err = w.wsconn.SetReadDeadline(t)
	if err != nil {
		return err
	}
	err = w.wsconn.SetWriteDeadline(t)
	if err != nil {
		return err
	}
	return nil
}

//
// SetReadDeadline
// @Description:
// @receiver w
// @param t
// @return error
//
func (w *WSConn) SetReadDeadline(t time.Time) error {
	return w.wsconn.SetReadDeadline(t)
}

//
// SetWriteDeadline
// @Description:
// @receiver w
// @param t
// @return error
//
func (w *WSConn) SetWriteDeadline(t time.Time) error {
	return w.wsconn.SetWriteDeadline(t)
}
