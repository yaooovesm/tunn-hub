package traffic

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"strings"
	"sync/atomic"
	"time"
)

//
// FlowStatisticsFP
// @Description:
//
type FlowStatisticsFP struct {
	Name            string
	Flow            uint64
	FlowStamp       uint64
	FlowCommitStamp uint64
	FlowSpeed       uint64
	Packet          uint64
	PacketStamp     uint64
	PacketSpeed     uint64
	RecordGap       int
	Print           bool
}

//
// Init
// @Description:
// @receiver fp
//
func (fp *FlowStatisticsFP) Init() bool {
	fp.Flow = 0
	fp.FlowStamp = 0
	fp.FlowSpeed = 0
	fp.FlowCommitStamp = 0
	fp.Packet = 0
	fp.PacketStamp = 0
	fp.PacketSpeed = 0
	if fp.RecordGap == 0 {
		fp.RecordGap = 1000
	}
	go fp.statistics()
	return true
}

//
// LoadFromDump
// @Description:
// @receiver fp
// @param file
// @return error
//
func (fp *FlowStatisticsFP) LoadFromDump(file string) error {
	bytes, err := ioutil.ReadFile(file)
	if err != nil && strings.Contains(err.Error(), "no such file or directory") || len(bytes) <= 0 {
		return errors.New("no content")
	}
	if err != nil {
		return err
	}
	decodeString, err := base64.StdEncoding.DecodeString(string(bytes))
	if err != nil {
		return err
	}
	temp := FlowStatisticsFP{}
	err = json.Unmarshal(decodeString, &temp)
	if err != nil {
		return err
	}
	fp.Flow = temp.Flow
	fp.FlowStamp = fp.Flow
	fp.FlowCommitStamp = fp.Flow
	fp.Packet = temp.Packet
	fp.PacketStamp = temp.PacketStamp
	return nil
}

//
// Dump
// @Description:
// @receiver fp
// @param file
// @return error
//
func (fp *FlowStatisticsFP) Dump(file string) error {
	marshal, err := json.Marshal(fp)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(file, []byte(base64.StdEncoding.EncodeToString(marshal)), 0600)
}

//
// Commit
// @Description:
// @receiver fp
// @return uint64
//
func (fp *FlowStatisticsFP) Commit() uint64 {
	size := fp.Flow - fp.FlowCommitStamp
	fp.FlowCommitStamp = fp.Flow
	return size
}

//
// Process
// @Description:
// @receiver fp
// @param raw
// @return []byte
//
func (fp *FlowStatisticsFP) Process(raw []byte) []byte {
	atomic.AddUint64(&fp.Packet, 1)
	atomic.AddUint64(&fp.Flow, uint64(len(raw)))
	//fp.Packet++
	//fp.Flow += uint64(len(raw))
	return raw
}

//
// calcPacket
// @Description:
// @receiver fp
//
func (fp *FlowStatisticsFP) calcPacket() {
	if fp.Packet == 0 {
		fp.PacketSpeed = 0
		return
	}
	fp.PacketSpeed = fp.Packet - fp.PacketStamp
	fp.PacketStamp = fp.Packet
}

//
// calcFlow
// @Description:
// @receiver fp
//
func (fp *FlowStatisticsFP) calcFlow() {
	if fp.Flow == 0 {
		fp.FlowSpeed = 0
		return
	}
	fp.FlowSpeed = fp.Flow - fp.FlowStamp
	fp.FlowStamp = fp.Flow
}

//
// Statistics
// @Description:
// @receiver fp
//
func (fp *FlowStatisticsFP) statistics() {
	for {
		time.Sleep(time.Millisecond * time.Duration(fp.RecordGap))
		fp.calcFlow()
		fp.calcPacket()
		if fp.Print {
			fmt.Println("[", fp.Name, "] packet_speed=", fp.PacketSpeed, "p/s rx_flow_speed=", fp.FlowSpeed/1024/1024, "mb/s")
		}
	}
}
