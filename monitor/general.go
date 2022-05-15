package monitor

//
// SystemData
// @Description:
//
type SystemData struct {
	CPU    CpuData    `json:"cpu"`
	Memory MemoryData `json:"memory"`
	Disk   DiskData   `json:"disk"`
}

//
// CollectAll
// @Description:
// @receiver d
//
func (d *SystemData) CollectAll() {
	_ = d.CPU.Collect()
	_ = d.Disk.Collect()
	_ = d.Memory.Collect()
}
