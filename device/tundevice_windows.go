package device

import (
	"errors"
	log "github.com/cihub/seelog"
	"golang.zx2c4.com/wireguard/tun"
	"net"
	"os/exec"
	"strconv"
	"tunn-hub/config"
	"tunn-hub/networking"
)

//
// TunDevice
// @Description:
//
type TunDevice struct {
	iface         tun.Device
	config        config.Device
	cidr          string
	mtu           int
	clearCIDRFunc func()
}

//
// Name
// @Description:
// @receiver d
// @return string
//
func (d *TunDevice) Name() string {
	if d.iface == nil {
		return ""
	}
	name, err := d.iface.Name()
	if err != nil {
		return ""
	}
	return name
}

//
// Create
// @Description:
// @receiver d
// @param config
// @return error
//
func (d *TunDevice) Create(config config.Config) error {
	dev, err := tun.CreateTUN(DefaultTunDeviceName, config.Global.MTU)
	if err != nil {
		return err
	}
	if config.Device.CIDR == "" {
		return errors.New("cidr not set")
	}
	d.config = config.Device
	d.mtu = config.Global.MTU
	d.iface = dev
	return nil
}

//
// Setup
// @Description:
// @receiver d
// @return error
//
func (d *TunDevice) Setup() error {
	err := d.setCIDR(d.config.CIDR)
	if err != nil {
		return err
	}
	err = d.setMTU(d.mtu)
	if err != nil {
		return err
	}
	routes := config.Current.Routes
	for i := range routes {
		if routes[i].Option == config.RouteOptionImport {
			log.Info("import route : ", routes[i].Network)
			networking.AddSystemRoute(routes[i].Network, d.Name())
		}
	}
	//auto up in windows
	return nil
}

//
// setCIDR
// @Description:
// @receiver d
// @param cidr
// @return error
//
func (d *TunDevice) setCIDR(cidr string) error {
	if d.iface == nil {
		return errors.New("device not found")
	}
	ip, ipNet, err := net.ParseCIDR(cidr)
	if err != nil {
		return err
	}
	name := d.Name()
	cmd := exec.Command("PowerShell",
		"netsh", "interface", "ip", "set", "Address", "Name=\""+name+"\"", "source=static", "addr="+ip.String(), "mask="+ipv4MaskString(ipNet.Mask), "store=active")
	err = cmd.Run()
	if err == nil {
		d.clearCIDRFunc = func() {
			cmd := exec.Command("PowerShell",
				"netsh", "interface", "ip", "delete", "Address", "Name=\""+name+"\"", "addr="+ip.String())
			_ = cmd.Run()
		}
		d.cidr = cidr
	}
	return err
}

//
// setMTU
// @Description:
// @receiver d
// @param mtu
// @return error
//
func (d *TunDevice) setMTU(mtu int) error {
	if d.iface == nil {
		return errors.New("device not found")
	}
	name := d.Name()
	cmd := exec.Command("PowerShell",
		"netsh", "interface", "ipv4", "set", "interface", "\""+name+"\"", "mtu="+strconv.Itoa(mtu))
	return cmd.Run()
}

//
// Close
// @Description:
// @receiver d
// @return error
//
func (d *TunDevice) Close() error {
	if d.iface == nil {
		return errors.New("device not found")
	}
	return d.iface.Close()
}

//
// Up
// @Description:
// @receiver d
// @return error
//
func (d *TunDevice) Up() error {
	if d.iface == nil {
		return errors.New("device not found")
	}
	cmd := exec.Command("PowerShell", "netsh", "interface", "set", "interface", d.Name(), "enabled")
	return cmd.Run()
}

//
// Down
// @Description:
// @receiver d
// @return error
//
func (d *TunDevice) Down() error {
	if d.iface == nil {
		return errors.New("device not found")
	}
	cmd := exec.Command("PowerShell", "netsh", "interface", "set", "interface", d.Name(), "disabled")
	return cmd.Run()
}

//
// Read
// @Description:
// @receiver d
// @param packet
// @return n
// @return err
//
func (d *TunDevice) Read(packet []byte) (n int, err error) {
	n, err = d.iface.Read(packet, 0)
	return
}

//
// Write
// @Description:
// @receiver d
// @param packet
// @return n
// @return err
//
func (d *TunDevice) Write(packet []byte) (n int, err error) {
	n, err = d.iface.Write(packet, 0)
	return
}

//
// OverwriteCIDR
// @Description:
// @receiver d
// @param cidr
// @return error
//
func (d *TunDevice) OverwriteCIDR(cidr string) error {
	if cidr == d.cidr {
		return nil
	}
	if d.clearCIDRFunc != nil {
		d.clearCIDRFunc()
	}
	return d.setCIDR(cidr)
}
