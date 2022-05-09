package device

import (
	"tunn-hub/common/config"
)

const DefaultTunDeviceName = "tunnel"

//
// Device
// @Description:
//
type Device interface {
	Name() string                           //return device name if device exists
	Create(config config.Config) error      //create device
	Close() error                           //close device
	Setup() error                           //setup device
	OverwriteCIDR(cidr string) error        //overwrite cidr
	Up() error                              //set device up
	Down() error                            //set device down
	Read(packet []byte) (n int, err error)  //read
	Write(packet []byte) (n int, err error) //write
}

//
// NewTunDevice
// @Description:
// @return device
// @return err
//
func NewTunDevice() (device Device, err error) {
	d := TunDevice{}
	err = d.Create(config.Current)
	return &d, err
}

//
// NewTunDeviceWithConfig
// @Description:
// @param cfg
// @return device
// @return err
//
func NewTunDeviceWithConfig(cfg config.Config) (device Device, err error) {
	d := TunDevice{}
	err = d.Create(cfg)
	return &d, err
}
