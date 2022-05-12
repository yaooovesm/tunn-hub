package transmitter

import (
	log "github.com/cihub/seelog"
	"net"
	"sync"
	"tunn-hub/traffic"
)

//
// MultiConn
// @Description:
//
type MultiConn struct {
	Current            *Tunnel
	conns              []*Tunnel
	ver                Version
	maxIndex           int
	currentIndex       int
	Size               int
	Name               string
	electionGap        int
	electionCount      int
	writeFlowProcessor *traffic.FlowProcessors
	sync.RWMutex
}

//
// NewMultiConn
// @Description:
// @param name
// @return *MultiConn
//
func NewMultiConn(name string) *MultiConn {
	m := &MultiConn{Name: name, electionGap: 20, ver: V2}
	return m
}

//
// SetWriteFlowProcessors
// @Description:
// @receiver m
// @param processors
//
func (m *MultiConn) SetWriteFlowProcessors(processors *traffic.FlowProcessors) {
	m.writeFlowProcessor = processors
}

//
// Push
// @Description:
// @receiver m
// @param conn
//
func (m *MultiConn) Push(conn net.Conn) int {
	m.Lock()
	num := m.Size
	tunnel := NewTunnel(conn, m.ver)
	if m.Current == nil {
		m.Current = tunnel
	}
	m.conns = append(m.conns, tunnel)
	m.Size++
	m.maxIndex = m.Size - 1
	m.Unlock()
	return num
	//log.Info("[multiConn:", m.Name, "] pushed connection, current size=", m.Size)
}

//
// Get
// @Description:
// @receiver m
//
func (m *MultiConn) Get() *Tunnel {
	defer func() {
		if m.Size <= 1 {
			return
		}
		m.electionCount++
		if m.electionCount == m.electionGap {
			m.electionCount = 0
			m.election()
		}
	}()
	return m.Current
}

//
// Write
// @Description:
// @receiver m
// @param pl
// @return n
// @return err
//
func (m *MultiConn) Write(pl []byte) (n int, err error) {
	if m.writeFlowProcessor != nil {
		pl = m.writeFlowProcessor.Process(pl)
	}
	n, err = m.Get().Write(pl)
	return
}

//
// election
// @Description:
// @receiver m
//
func (m *MultiConn) election() {
	m.Lock()
	if m.currentIndex < m.maxIndex {
		m.currentIndex++
	} else {
		m.currentIndex = 0
	}
	m.Current = m.conns[m.currentIndex]
	m.Unlock()
	//log.Info("[multiConn:", m.Name, "] current connection : ", m.currentIndex)
}

//
// Close
// @Description:
// @receiver m
//
func (m *MultiConn) Close() {
	m.Lock()
	for i := range m.conns {
		if m.conns[i] != nil {
			_ = m.conns[i].Close()
		}
	}
	m.Size = 0
	log.Info("[multiConn:", m.Name, "] close")
	m.Unlock()
}
