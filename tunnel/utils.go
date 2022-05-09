package tunnel

import (
	log "github.com/cihub/seelog"
	"net"
	"os/exec"
	"runtime"
	"strings"
)

//
// command
// @Description:
// @param c
// @param args
//
func command(c string, args ...string) {
	cmd := exec.Command(c, args...)
	err := cmd.Run()
	if err != nil {
		_ = log.Warn(c, " exec failed : ", err.Error())
	}
}

//
// addRouteLinux
// @Description:
// @param network
// @param dev
//
func addRouteLinux(network string, dev string) {
	var cmd *exec.Cmd
	if runtime.GOOS == "linux" {
		cmd = exec.Command("/sbin/ip", "route", "add", network, "dev", dev)
	} else if runtime.GOOS == "windows" {

	} else {
		log.Info("os not support : ", runtime.GOOS)
		return
	}
	err := cmd.Run()
	if err != nil {
		log.Info("[", dev, "] import route failed <-- ", network, " : ", err.Error())
		return
	}
	log.Info("[", dev, "] import route <-- ", network)
}

//
// GetMacByInetName
// @Description:
// @param Name
//
func GetMacByInetName(name string) string {
	interfaces, err := net.Interfaces()
	if err != nil {
		return ""
	}
	for _, inter := range interfaces {
		if inter.Name == name {
			return strings.ToLower(inter.HardwareAddr.String())
		}
	}
	return ""
}

//
// HasSameIp
// @Description:
// @param Name
//
func HasSameIp(ip string) bool {
	interfaces, err := net.Interfaces()
	if err != nil {
		return true
	}
	for _, inter := range interfaces {
		if addrs, err := inter.Addrs(); err != nil {
			return true
		} else {
			for i := range addrs {
				if addrs[i].String() == ip {
					log.Info(inter.Name, " ip duplicate")
					return true
				}
			}
		}
	}
	return false
}
