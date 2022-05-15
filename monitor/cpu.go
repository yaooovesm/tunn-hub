package monitor

import (
	"github.com/shirou/gopsutil/v3/cpu"
	"time"
)

//
// CpuData
// @Description:
//
type CpuData struct {
	Usage float64 `json:"usage"`
	Error string  `json:"error"`
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
	return nil
}
