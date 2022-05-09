package networking

import (
	"fmt"
	"net"
	"testing"
)

func TestNewIPAddressPool(t *testing.T) {
	ipRange := IPRange{}
	ipRange.Start("192.0.1.250").End("192.0.2.2")
	m := ipRange.Map()
	_, ipNet, err := net.ParseCIDR("192.0.2.0/16")
	if err != nil {
		return
	}
	size := ipRange.Size(ipNet)
	fmt.Println(m)
	fmt.Println(len(m), " / ", size)
	pool, err := NewIPv4AddressPool("192.0.2.0/16", ipRange)
	if err != nil {
		fmt.Println("create : ", err)
		return
	}
	for i := 0; i < 7; i++ {
		dispatch, err := pool.DispatchCIDR("")
		if err != nil {
			fmt.Println("error dispatch@", i, " : ", err)
			continue
		}
		fmt.Println("dispatch ", i, " :", dispatch)
	}
	fmt.Println("used : ", pool.used)
	fmt.Println("return back : 192.0.1.253")
	pool.ReturnBack("192.0.1.253")
	fmt.Println("return back : 192.0.2.1")
	pool.ReturnBack("192.0.2.1")
	fmt.Println("used : ", pool.used)
	fmt.Println("pick 192.0.2.1")
	pool.PickCIDR("192.0.2.1/24", "")
	fmt.Println("used : ", pool.used)
	for i := 0; i < 3; i++ {
		dispatch, err := pool.DispatchCIDR("")
		if err != nil {
			fmt.Println("error dispatch@", i, " : ", err)
			continue
		}
		fmt.Println("dispatch ", i, " :", dispatch)
	}
	fmt.Println("used : ", pool.used)
}
