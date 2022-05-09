package model

import "tunn-hub/traffic"

//
// LinkStatus
// @Description:
//
type LinkStatus struct {
	Flow        uint64
	FlowSpeed   uint64
	Packet      uint64
	PacketSpeed uint64
}

//
// ReadFromFP
// @Description:
// @receiver ls
// @param fp
//
func (ls *LinkStatus) ReadFromFP(fp *traffic.FlowStatisticsFP) {
	if fp == nil {
		return
	}
	ls.Flow = fp.Flow
	ls.Packet = fp.Packet
	ls.FlowSpeed = fp.FlowSpeed
	ls.PacketSpeed = fp.PacketSpeed
}
