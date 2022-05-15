package monitor

import (
	"github.com/shirou/gopsutil/v3/disk"
)

//
// DiskData
// @Description:
//
type DiskData struct {
	Total uint64  `json:"total"`
	Used  uint64  `json:"used"`
	Usage float64 `json:"usage"`
	Error string  `json:"error"`
}

//
// Collect
// @Description:
// @receiver d
// @return error
//
func (d *DiskData) Collect() error {
	parts, err := disk.Partitions(true)
	if err != nil {
		d.Error = err.Error()
		return err
	}
	total := uint64(0)
	used := uint64(0)
	for _, part := range parts {
		diskInfo, _ := disk.Usage(part.Mountpoint)
		total += diskInfo.Total
		used += diskInfo.Used
	}
	d.Used = used
	d.Total = total
	d.Usage = float64(used) / float64(total) * 100
	return nil
}
