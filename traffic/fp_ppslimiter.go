package traffic

import (
	"math"
	"time"
)

//
// LimitedPPSBucket
// @Description:
//
type LimitedPPSBucket struct {
	bk      chan byte //bucket
	gap     int       //流量间隙
	mtu     int       //最大传输单元
	maxpps  int       //最大pps
	currpps int       //当前pps
	latency int       //生产间隔
	running bool      //是否运行中
}

//
// NewLimitedPPSBucket
// @Description:
// @return LimitedPPSBucket
//
func NewLimitedPPSBucket(maxSpeed uint64, mtu int) *LimitedPPSBucket {
	//计算最大pps
	pps := int(maxSpeed / uint64(mtu))
	if maxSpeed%uint64(mtu) > 0 {
		pps += 1
	}
	bucket := LimitedPPSBucket{
		bk:      make(chan byte, pps),
		maxpps:  pps,
		latency: 200,
		running: true,
		gap:     0,
		currpps: 0,
		mtu:     mtu,
	}
	go bucket.produce()
	return &bucket
}

//
// produce
// @Description:
// @receiver l
//
func (l *LimitedPPSBucket) produce() {
	//动态限速
	//速度为每秒生l.speed个包,每个包大小为l.mtu,总大小为l.limit
	//将每秒分成1000/l.latency份,每份延迟l.latency执行，否则会出现较大的速度波动

	slice := 1000 / l.latency
	unit := l.maxpps / slice
	//扩容方法
	fill := func(count int) {
		for i := 0; i < count; i++ {
			l.bk <- 1
		}
	}
	var max int
	for l.running {
		//每次生产的数额<=unit
		pps := l.currpps * slice
		bk := len(l.bk)
		max = unit * ((pps / unit) + 1)
		if max > l.maxpps {
			max = l.maxpps
		}
		if bk < max {
			diff := int(math.Abs(float64(bk - unit)))
			if diff > unit {
				fill(unit)
			} else {
				fill(diff)
			}
		}
		//每个周期定时清空数据
		l.currpps = 0
		time.Sleep(time.Millisecond * time.Duration(l.latency))
	}
}

//
// Take
// @Description: 从桶中取令牌，如果余量不足以取出则等待
// @receiver l
//
func (l *LimitedPPSBucket) Take(p []byte) {
	if l.running {
		//当积累的间隔足够多时额外补充数据包
		if l.gap >= l.mtu {
			l.gap -= l.mtu
		} else {
			//等待生产包
			<-l.bk
		}
		//补充gap
		length := len(p)
		l.gap += l.mtu - length
		//l.fpl += uint64(length)
		l.currpps += 1
	}
}

//
// Close
// @Description:
// @receiver l
//
func (l *LimitedPPSBucket) Close() {
	defer func() {
		recover()
	}()
	//有可能已关闭
	close(l.bk)
	l.running = false
}

//
// PPSLimiter
// @Description:
//
type PPSLimiter struct {
	bandwidth int
	mtu       int
	bucket    *LimitedPPSBucket
}

//
// NewPPSLimiterFP
// @Description:
// @param bandwidth
// @param mtu
// @return *PPSLimiter
//
func NewPPSLimiterFP(bandwidth int, mtu int) PPSLimiter {
	return PPSLimiter{
		bandwidth: bandwidth,
		mtu:       mtu,
	}
}

//
// Init
// @Description:
// @receiver PPSLimiter
// @return bool
//
func (l *PPSLimiter) Init() bool {
	if l.bandwidth <= 0 {
		return false
	}
	l.bucket = NewLimitedPPSBucket(uint64(l.bandwidth*1024*1024/8), l.mtu)
	return true
}

//
// Process
// @Description:
// @receiver l
// @param raw
// @return []byte
//
func (l *PPSLimiter) Process(raw []byte) []byte {
	l.bucket.Take(raw)
	return raw
}

//
// Close
// @Description:
// @receiver l
//
func (l *PPSLimiter) Close() {
	l.bucket.Close()
}
