package traffic

import (
	"time"
)

//
// FlowLimitFP
// @Description:
//
type FlowLimitFP struct {
	counter  *FlowStatisticsFP //计数器
	Flow     uint64            //起始点
	Limit    uint64            //限制
	kickFunc func()            //断开方法
	running  bool
}

//
// NewFlowLimitFP
// @Description:
// @return *FlowLimitFP
//
func NewFlowLimitFP(counter *FlowStatisticsFP, flow, limit uint64, kick func()) FlowLimitFP {
	return FlowLimitFP{
		counter:  counter,
		Flow:     flow,
		Limit:    limit,
		kickFunc: kick,
		running:  false,
	}
}

//
// Init
// @Description:
// @receiver f
// @return bool
//
func (f *FlowLimitFP) Init() bool {
	f.running = true
	go func() {
		for f.running {
			time.Sleep(time.Second * 20)
			if f.Limit != 0 && (f.counter.Flow+f.Flow) >= f.Limit {
				f.kickFunc()
			}
		}
	}()
	return true
}

//
// Process
// @Description:
// @receiver f
// @param raw
// @return []byte
//
func (f *FlowLimitFP) Process(raw []byte) []byte {
	return raw
}

//
// Close
// @Description:
// @receiver f
//
func (f *FlowLimitFP) Close() {
	f.running = false
}
