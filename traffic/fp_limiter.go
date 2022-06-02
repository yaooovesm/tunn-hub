package traffic

import (
	"time"
)

//
// LimitedBucket
// @Description:
//
type LimitedBucket struct {
	limit   uint64
	mtu     int
	bk      chan byte
	pps     int
	latency int
	running bool
}

//
// NewLimitedBucket
// @Description:
// @return LimitedBucket
//
func NewLimitedBucket(maxSpeed uint64, mtu int) *LimitedBucket {
	//计算最大pps
	pps := int(maxSpeed / uint64(mtu))
	if maxSpeed%uint64(mtu) > 0 {
		pps += 1
	}
	bucket := &LimitedBucket{
		bk:      make(chan byte, pps),
		limit:   maxSpeed,
		mtu:     mtu,
		pps:     pps,
		latency: 200,
		running: true,
	}
	go bucket.produce()
	return bucket
}

//
// produce
// @Description:
// @receiver l
//
func (l *LimitedBucket) produce() {
	//以恒定速率生产
	//速度为每秒生l.speed个包,每个包大小为l.mtu,总大小为l.limit
	//将每秒分成1000/l.latency份,每份延迟l.latency执行，否则会出现较大的速度波动
	slice := 1000 / l.latency
	unit := l.pps / slice
	for l.running {
		for i := 0; i < unit; i++ {
			l.bk <- 1
		}
		time.Sleep(time.Millisecond * time.Duration(l.latency))
	}
}

//
// Take
// @Description: 从桶中取令牌，如果余量不足以取出则等待
// @receiver l
//
func (l *LimitedBucket) Take() {
	if l.running {
		<-l.bk
	}
}

//
// Close
// @Description:
// @receiver l
//
func (l *LimitedBucket) Close() {
	l.running = false
}

//
// Limiter
// @Description:
//
type Limiter struct {
	bandwidth int
	mtu       int
	bucket    *LimitedBucket
}

//
// NewLimiterFP
// @Description:
// @param bandwidth
// @param mtu
// @return *Limiter
//
func NewLimiterFP(bandwidth int, mtu int) Limiter {
	return Limiter{
		bandwidth: bandwidth,
		mtu:       mtu,
	}
}

//
// Init
// @Description:
// @receiver Limiter
// @return bool
//
func (l *Limiter) Init() bool {
	if l.bandwidth <= 0 {
		return false
	}
	l.bucket = NewLimitedBucket(uint64(l.bandwidth*1024*1024/8), l.mtu)
	return true
}

//
// Process
// @Description:
// @receiver l
// @param raw
// @return []byte
//
func (l *Limiter) Process(raw []byte) []byte {
	l.bucket.Take()
	return raw
}

//
// Close
// @Description:
// @receiver l
//
func (l *Limiter) Close() {
	l.bucket.Close()
}
