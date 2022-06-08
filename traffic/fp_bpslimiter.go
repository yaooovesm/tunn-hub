package traffic

import "time"

type LimitedBPSBucket struct {
	token          uint64
	max            uint64
	mtu            uint64
	produceLatency int
	unitLatency    time.Duration
	running        bool
}

func NewLimitedBPSBucket(maxSpeed uint64, mtu int) *LimitedBPSBucket {
	bucket := &LimitedBPSBucket{
		max:            maxSpeed,
		produceLatency: 200,
		unitLatency:    time.Duration((float64(mtu) / float64(maxSpeed)) * 1000000),
		mtu:            uint64(mtu),
		running:        true,
	}
	go bucket.produce()
	return bucket
}

//
// produce
// @Description:
// @receiver l
//
func (l *LimitedBPSBucket) produce() {
	u := (l.max * uint64(l.produceLatency)) / 1000
	max := l.max - u
	half := l.max / 2
	for l.running {
		if l.token <= max && l.token < half {
			l.token += u
		}
		time.Sleep(time.Millisecond * time.Duration(l.produceLatency))
	}
}

//
// Take
// @Description:
// @receiver l
//
func (l *LimitedBPSBucket) Take(p []byte) {
	required := uint64(len(p))
retry:
	if l.token >= required {
		l.token -= required
	} else {
		time.Sleep(time.Nanosecond * l.unitLatency)
		goto retry
	}
}

//
// Close
// @Description:
// @receiver l
//
func (l *LimitedBPSBucket) Close() {
	l.running = false
}

//
// BPSLimiter
// @Description:
//
type BPSLimiter struct {
	bandwidth int
	mtu       int
	bucket    *LimitedBPSBucket
}

//
// NewBPSLimiterFP
// @Description:
// @param bandwidth
// @param mtu
// @return BPSLimiter
//
func NewBPSLimiterFP(bandwidth int, mtu int) BPSLimiter {
	return BPSLimiter{
		bandwidth: bandwidth,
		mtu:       mtu,
	}
}

//
// Init
// @Description:
// @receiver l
// @return bool
//
func (l *BPSLimiter) Init() bool {
	if l.bandwidth <= 0 {
		return false
	}
	l.bucket = NewLimitedBPSBucket(uint64(l.bandwidth*1024*1024/8), l.mtu)
	return true
}

//
// Process
// @Description:
// @receiver l
// @param raw
// @return []byte
//
func (l *BPSLimiter) Process(raw []byte) []byte {
	l.bucket.Take(raw)
	return raw
}

//
// Close
// @Description:
// @receiver l
//
func (l *BPSLimiter) Close() {
	l.bucket.Close()
}
