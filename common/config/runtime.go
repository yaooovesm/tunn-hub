package config

import (
	"github.com/shirou/gopsutil/v3/host"
	"runtime"
	"tunn-hub/version"
)

//
// Runtime
// @Description:
//
type Runtime struct {
	OS       string `json:"os"`
	Version  string `json:"version"`
	Arch     string `json:"arch"`
	Platform string `json:"platform"`
	App      string `json:"app"`
}

//
// Collect
// @Description:
// @receiver r
//
func (r *Runtime) Collect() {
	info, err := host.Info()
	r.App = version.Version
	if err != nil {
		r.OS = runtime.GOOS
		r.Arch = runtime.GOARCH
		r.Platform = "unknown"
		r.Version = "unknown"
		return
	}
	r.OS = info.OS
	r.Platform = info.Platform
	r.Arch = info.KernelArch
	r.Version = info.PlatformVersion
}
