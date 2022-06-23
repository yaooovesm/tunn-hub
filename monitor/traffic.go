package monitor

import (
	"encoding/json"
	"io/ioutil"
	"time"
	"tunn-hub/traffic"
)

//
// HubTrafficRecorder
// @Description:
//
type HubTrafficRecorder struct {
	old  []HubTrafficStamp
	data []HubTrafficStamp
}

//
// NewHubTrafficRecorder
// @Description:
// @return HubTrafficRecorder
//
func NewHubTrafficRecorder() HubTrafficRecorder {
	return HubTrafficRecorder{data: make([]HubTrafficStamp, 0)}
}

//
// Recent
// @Description:
// @receiver r
// @param count
// @return []HubTrafficStamp
//
func (r *HubTrafficRecorder) Recent(count int) []HubTrafficStamp {
	arr := make([]HubTrafficStamp, 0)
	currLen := len(r.data)
	if currLen < count {
		tmp := append(r.old, r.data...)
		require := len(tmp) - count
		if require < 0 {
			require = 0
		}
		arr = tmp[require:]
	} else {
		arr = append(arr, r.data[currLen-count:]...)
	}
	return arr
}

//
// Latest
// @Description:
// @receiver r
// @return []HubTrafficStamp
//
func (r *HubTrafficRecorder) Latest() []HubTrafficStamp {
	return r.data
}

//
// Push
// @Description:
// @receiver r
// @param stamp
//
func (r *HubTrafficRecorder) Push(stamp HubTrafficStamp) {
	r.data = append(r.data, stamp)
}

//
// DumpAndReset
// @Description:
// @receiver r
// @param file
// @return error
//
func (r *HubTrafficRecorder) DumpAndReset(file string) error {
	bytes, err := json.Marshal(r.data)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(file, bytes, 0644)
	if err != nil {
		return err
	}
	r.old = make([]HubTrafficStamp, len(r.data))
	copy(r.old, r.data)
	r.data = make([]HubTrafficStamp, 0)
	return nil
}

//
// HubTrafficStamp
// @Description:
//
type HubTrafficStamp struct {
	RxFlow    uint64 `json:"rx_flow"`   //接收流量
	TxFlow    uint64 `json:"tx_flow"`   //发送流量
	RxPacket  uint64 `json:"rx_packet"` //数据包接收
	TxPacket  uint64 `json:"tx_packet"` //数据包发送
	Timestamp int64  `json:"timestamp"` //时间戳
}

//
// CreateHubTrafficStampWithDuration
// @Description:
// @param rx
// @param tx
// @param gap
// @param slice
// @return HubTrafficStamp
//
func CreateHubTrafficStampWithDuration(rx, tx *traffic.FlowStatisticsFP, gap int) HubTrafficStamp {
	if rx == nil || tx == nil {
		return HubTrafficStamp{}
	}
	start := time.Now().UnixMilli()
	rxFlowSt1 := rx.Flow
	txFlowSt1 := tx.Flow
	rxPacketSt1 := rx.Packet
	txPacketSt1 := tx.Packet
	time.Sleep(time.Millisecond * time.Duration(gap))
	rxFlowSt2 := rx.Flow
	txFlowSt2 := tx.Flow
	rxPacketSt2 := rx.Packet
	txPacketSt2 := tx.Packet
	return HubTrafficStamp{
		RxFlow:    rxFlowSt2 - rxFlowSt1,
		TxFlow:    txFlowSt2 - txFlowSt1,
		RxPacket:  rxPacketSt2 - rxPacketSt1,
		TxPacket:  txPacketSt2 - txPacketSt1,
		Timestamp: start,
	}
}
