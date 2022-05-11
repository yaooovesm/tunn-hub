package device

import (
	"errors"
	"github.com/songgao/water"
	"os/exec"
	"strconv"
	"tunn-hub/config"
)

//
// TunDevice
// @Description:
//
type TunDevice struct {
	iface         *water.Interface
	config        config.Device
	mtu           int
	cidr          string
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
	return d.iface.Name()
}

//
// Create
// @Description:
// @receiver d
// @param config
// @return error
//
func (d *TunDevice) Create(config config.Config) error {
	dev, err := water.New(water.Config{
		DeviceType: water.TUN,
		PlatformSpecificParams: water.PlatformSpecificParams{
			Name: DefaultTunDeviceName,
		},
	})
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
	err = d.setMtu()
	if err != nil {
		return err
	}
	err = d.Up()
	if err != nil {
		return err
	}
	//不在此处引入路由，由systemrt.SystemRouteTable统一托管
	//routes := config.Current.Routes
	//for i := range routes {
	//	if routes[i].Option == config.RouteOptionImport {
	//		log.Info("import route : ", routes[i].Network)
	//		networking.AddSystemRoute(routes[i].Network, d.Name())
	//	}
	//}
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
	name := d.Name()
	cmd := exec.Command("/sbin/ip", "address", "add", cidr, "dev", name)
	err := cmd.Run()
	if err == nil {
		d.clearCIDRFunc = func() {
			cmd := exec.Command("/sbin/ip", "address", "del", cidr, "dev", name)
			_ = cmd.Run()
		}
		d.cidr = cidr
	}
	return nil
}

//
// setMtu
// @Description:
// @receiver d
// @return error
//
func (d *TunDevice) setMtu() error {
	if d.iface == nil {
		return errors.New("device not found")
	}
	name := d.Name()
	cmd := exec.Command("/sbin/ip", "link", "set", "dev", name, "mtu", strconv.Itoa(d.mtu))
	_ = cmd.Run()
	return nil
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
	cmd := exec.Command("/sbin/ip", "link", "set", "dev", d.Name(), "up")
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
	cmd := exec.Command("/sbin/ip", "link", "set", "dev", d.Name(), "down")
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
	n, err = d.iface.Read(packet)
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
	n, err = d.iface.Write(packet)
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
	err := d.setCIDR(cidr)
	if err != nil {
		return err
	}
	return nil
}
