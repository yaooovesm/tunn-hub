package networking

import (
	"errors"
	"fmt"
	log "github.com/cihub/seelog"
	"net"
	"os/exec"
	"strconv"
)

//
// AddSystemRoute
// @Description:
// @param network
// @param dev
//
func AddSystemRoute(network string, dev string) {
	log.Info("[", dev, "]add system route : ", network)
	ip, ipNet, err := net.ParseCIDR(network)
	if err != nil {
		_ = log.Warn("import ", network, " failed : ", err)
		return
	}
	//PowerShell route add -p [network] mask [mask] [dev_ip]
	devIp, index, err := getIpv4ByInterfaceName(dev)
	if err != nil {
		_ = log.Warn("import ", network, " failed : ", err)
		return
	}
	err = command("PowerShell", "route", "add", ip.String(), "mask", ipv4MaskString(ipNet.Mask), devIp, "IF", strconv.Itoa(index))
	if err != nil {
		_ = log.Warn("import ", network, " failed : ", err)
		return
	}
}

//
// command
// @Description:
// @param c
// @param args
//
func command(c string, args ...string) error {
	cmd := exec.Command(c, args...)
	log.Info("exec : ", cmd.String())
	return cmd.Run()
}

//
// ipv4MaskString
// @Description:
// @param m
// @return string
//
func ipv4MaskString(m []byte) string {
	if len(m) != 4 {
		panic("ipv4Mask: len must be 4 bytes")
	}
	return fmt.Sprintf("%d.%d.%d.%d", m[0], m[1], m[2], m[3])
}

//
// getIpv4ByInterfaceName
// @Description:
// @param name
// @return string
// @return error
//
func getIpv4ByInterfaceName(name string) (string, int, error) {
	interfaces, err := net.Interfaces()
	if err != nil {
		return "", -1, err
	}
	for _, inter := range interfaces {
		if inter.Name == name {
			addrs, err := inter.Addrs()
			if err != nil {
				return "", -1, err
			}
			for i := range addrs {
				if isv4(addrs[i]) {
					ip, _, err := net.ParseCIDR(addrs[i].String())
					if err != nil {
						return "", -1, err
					}
					return ip.String(), inter.Index, nil
				}
			}
		}
	}
	return "", -1, errors.New("interface not found")
}

//
// isv4
// @Description:
// @param addr
// @return bool
//
func isv4(addr net.Addr) bool {
	ip := addr.String()
	for i := 0; i < len(ip); i++ {
		switch ip[i] {
		case '.':
			return true
		case ':':
			return false
		}
	}
	return false
}
