package monitor

import (
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/process"
	"os"
	"time"
)

//
// CpuData
// @Description:
//
type CpuData struct {
	Usage   float64 `json:"usage"`
	AppUsed float64 `json:"app_used"`
	Error   string  `json:"error"`
}

//
// Collect
// @Description:
// @receiver c
// @return error
//
func (c *CpuData) Collect() error {
	percent, err := cpu.Percent(time.Second, false)
	if err != nil {
		c.Error = err.Error()
		return err
	}
	c.Usage = percent[0]
	p, err := process.NewProcess(int32(os.Getpid()))
	if err != nil {
		c.AppUsed = 0
		c.Error = err.Error()
	}
	cpuPercent, err := p.CPUPercent()
	if err != nil {
		c.AppUsed = 0
		c.Error = err.Error()
	} else {
		if cpuPercent > c.Usage {
			c.AppUsed = c.Usage
		} else {
			c.AppUsed = cpuPercent
		}
	}
	return nil
}
