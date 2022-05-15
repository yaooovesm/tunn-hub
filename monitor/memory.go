package monitor

import (
	"github.com/shirou/gopsutil/v3/mem"
	"github.com/shirou/gopsutil/v3/process"
	"os"
)

//
// MemoryData
// @Description:
//
type MemoryData struct {
	Total     uint64  `json:"total"`
	Used      uint64  `json:"used"`
	Usage     float64 `json:"usage"`
	AppUsed   uint64  `json:"app_used"`
	SwapTotal uint64  `json:"swap_total"`
	SwapUsed  uint64  `json:"swap_used"`
	SwapUsage float64 `json:"swap_usage"`
	Error     string  `json:"error"`
}

//
// Collect
// @Description:
// @receiver m
//
func (m *MemoryData) Collect() error {
	v, err := mem.VirtualMemory()
	if err != nil {
		m.Error = err.Error()
		return err
	}
	swapUsed := v.SwapTotal - v.SwapFree
	swapUsage := float64(0)
	if v.SwapTotal != 0 {
		swapUsage = float64((swapUsed) / v.SwapTotal)
	}
	m.Total = v.Total
	m.Used = v.Used
	m.Usage = v.UsedPercent
	m.SwapTotal = v.SwapTotal
	m.SwapUsed = swapUsed
	m.SwapUsage = swapUsage
	appMemoryUsage, err := collectAppMemoryUsage()
	if err != nil {
		m.AppUsed = 0
		m.Error = err.Error()
	} else {
		m.AppUsed = appMemoryUsage
	}
	//var ms runtime.MemStats
	//runtime.ReadMemStats(&ms)
	//m.AppUsed = ms.HeapIdle - ms.HeapReleased
	return nil
}

//
// collectAppMemoryUsage
// @Description:
// @return uint64
// @return error
//
func collectAppMemoryUsage() (uint64, error) {
	p, err := process.NewProcess(int32(os.Getpid()))
	if err != nil {
		return 0, err
	}
	info, err := p.MemoryInfo()
	if err != nil {
		return 0, err
	}
	return info.RSS, nil
}
