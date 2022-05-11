package networking

import (
	log "github.com/cihub/seelog"
	"os/exec"
)

//
// AddSystemRoute
// @Description:
// @param network
// @param dev
//
func AddSystemRoute(network string, dev string) error {
	log.Info("[", dev, "]add system route : ", network)
	err := command("/sbin/ip", "route", "add", network, "dev", dev)
	if err != nil {
		return log.Warn("import ", network, " failed : ", err)
	}
	return nil
}

//
// command
// @Description:
// @param c
// @param args
//
func command(c string, args ...string) error {
	cmd := exec.Command(c, args...)
	return cmd.Run()
}
