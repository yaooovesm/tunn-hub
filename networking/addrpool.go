package networking

import (
	"errors"
	"net"
	"strconv"
	"sync"
	"time"
)

//
// IPRange
// @Description:
//
type IPRange struct {
	start net.IP
	end   net.IP
}

//
// Start
// @Description:
// @receiver r
// @param ip
// @return *IPRange
//
func (r *IPRange) Start(ip string) *IPRange {
	r.start = net.ParseIP(ip)
	return r
}

//
// End
// @Description:
// @receiver r
// @param ip
// @return *IPRange
//
func (r *IPRange) End(ip string) *IPRange {
	r.end = net.ParseIP(ip)
	return r
}

//
// Size
// @Description:
// @receiver r
// @param subnet
// @return int
//
func (r IPRange) Size(subnet *net.IPNet) int {
	size, bits := subnet.Mask.Size()
	gap := 0
	if bits-size > 8 {
		gap = int(r.end.To4()[2] - r.start.To4()[2])
	}
	startHost := int(r.start.To4()[3])
	endHost := int(r.end.To4()[3])
	return 254*gap - startHost + endHost + 1
}

//
// Map
// @Description:
// @receiver r
// @return map[int]net.IP
//
func (r IPRange) Map() map[int]net.IP {
	m := map[int]net.IP{}
	var list []net.IP
	if r.end.To4()[2]-r.start.To4()[2] == 0 {
		for i := r.start.To4()[3]; i <= r.end.To4()[3]; i++ {
			list = append(list, []byte{r.start.To4()[0], r.start.To4()[1], r.start.To4()[2], i})
		}
	} else {
		//start
		for i := r.start.To4()[3]; i <= 254; i++ {
			list = append(list, []byte{r.start.To4()[0], r.start.To4()[1], r.start.To4()[2], i})
		}
		//middle
		for i := r.start.To4()[2] + 1; i < r.end.To4()[2]; i++ {
			for j := 1; j <= 254; j++ {
				list = append(list, []byte{r.end.To4()[0], r.end.To4()[1], i, byte(j)})
			}
		}
		//end
		for i := 1; i <= int(r.end.To4()[3]); i++ {
			list = append(list, []byte{r.end.To4()[0], r.end.To4()[1], r.end.To4()[2], byte(i)})
		}
	}
	for i := range list {
		m[i] = list[i]
	}
	return m
}

//
// IPAllocInfo
// @Description:
//
type IPAllocInfo struct {
	Date      int64
	Expire    int64
	UUID      string
	Address   string
	Network   string
	IsDynamic bool
}

//
// IPAddressPool
// @Description:
//
type IPAddressPool struct {
	ipNet       *net.IPNet
	ipRange     IPRange
	used        []net.IP
	ipTable     map[int]net.IP         //used_index:ip
	staticTable map[string]IPAllocInfo //ip:info
	info        map[string]IPAllocInfo //address:info
	size        int
	sync.RWMutex
}

//
// NewIPv4AddressPool
// @Description:
// @param network
// @param ipRange
// @return pool
// @return err
//
func NewIPv4AddressPool(network string, ipRange IPRange) (pool *IPAddressPool, err error) {
	ip, ipNet, err := net.ParseCIDR(network)
	if err != nil {
		return nil, err
	}
	if ip.To4() == nil {
		return nil, errors.New("invalid ipv4 network")
	}
	if !ipNet.Contains(ipRange.end) || !ipNet.Contains(ipRange.start) {
		return nil, errors.New("ip range out of network")
	}
	size := ipRange.Size(ipNet)
	return &IPAddressPool{
		ipNet:       ipNet,
		ipRange:     ipRange,
		used:        make([]net.IP, size),
		size:        size,
		ipTable:     ipRange.Map(),
		staticTable: map[string]IPAllocInfo{},
		info:        map[string]IPAllocInfo{},
	}, nil
}

//
// Info
// @Description:
// @receiver p
// @return map[string]IPAllocInfo
//
func (p *IPAddressPool) Info() map[string]IPAllocInfo {
	return p.info
}

//
// General
// @Description:
// @receiver p
// @return map[string]interface{}
//
func (p *IPAddressPool) General() map[string]interface{} {
	used := 0
	for i := range p.used {
		if p.used[i] != nil {
			used++
		}
	}
	mask, _ := p.ipNet.Mask.Size()
	return map[string]interface{}{
		"size":    p.size,
		"used":    used,
		"network": p.ipNet.IP.String() + "/" + strconv.Itoa(mask),
		"static":  len(p.staticTable),
		"start":   p.ipRange.start.To4().String(),
		"end":     p.ipRange.end.To4().String(),
	}
}

//
// allocInfo
// @Description:
// @receiver p
// @param uuid
// @param address
// @return IPAllocInfo
//
func (p *IPAddressPool) allocInfo(uuid string, address string, isDynamic bool) IPAllocInfo {
	mask, _ := p.ipNet.Mask.Size()
	return IPAllocInfo{
		Date:      time.Now().UnixMilli(),
		Expire:    0,
		UUID:      uuid,
		Address:   address,
		Network:   p.ipNet.IP.String() + "/" + strconv.Itoa(mask),
		IsDynamic: isDynamic,
	}
}

//
// available
// @Description:
// @receiver p
// @return bool
//
func (p *IPAddressPool) available() bool {
	used := 0
	for i := range p.used {
		if p.used[i] != nil {
			used++
		}
	}
	return p.ipRange.Size(p.ipNet)-used > 0
}

//
// DispatchCIDR
// @Description:
// @receiver p
// @return net.IP
//
func (p *IPAddressPool) DispatchCIDR(uuid string) (ip string, err error) {
	//??????????????????????????????IP
	p.Lock()
	defer p.Unlock()
	if !p.available() {
		return "", errors.New("no ip can be dispatch in pool")
	}
	//?????????????????????
	var index = 0
	for i := 0; i < len(p.used); i++ {
		if p.used[i] == nil {
			index = i
			break
		}
	}
	dispatch := p.ipTable[index]
	size, _ := p.ipNet.Mask.Size()
	cidr := dispatch.String() + "/" + strconv.Itoa(size)
	p.used[index] = dispatch
	p.info[dispatch.String()] = p.allocInfo(uuid, cidr, true)
	return cidr, nil
}

//
// StaticCIDR
// @Description:
// @receiver p
// @param uuid
// @param cidr
// @return string
// @return error
//
func (p *IPAddressPool) StaticCIDR(uuid string, cidr string) (string, error) {
	p.Lock()
	defer p.Unlock()
	ip, _, err := net.ParseCIDR(cidr)
	if err != nil {
		return "", err
	}
	//????????????????????????????????????
	for u := range p.staticTable {
		if p.staticTable[u].Address == ip.String() {
			return "", errors.New("static ip duplicate")
		}
	}
	//?????????
	//???????????????????????????
	info := p.allocInfo(uuid, cidr, false)
	for i := range p.ipTable {
		if p.ipTable[i].String() == ip.String() {
			//???????????????
			if p.used[i] != nil {
				//????????????
				return "", errors.New("address is already in use")
			} else {
				//????????????
				p.staticTable[ip.String()] = info
				p.used[i] = p.ipTable[i]
				p.info[ip.String()] = info
				return cidr, nil
			}
		}
	}
	//??????????????????
	p.staticTable[ip.String()] = info
	p.info[ip.String()] = info
	return cidr, nil
}

//
// ReturnBack
// @Description:
// @receiver p
// @param ip
//
func (p *IPAddressPool) ReturnBack(ip string) {
	p.Lock()
	defer p.Unlock()
	for i := range p.used {
		if p.used[i].String() == ip {
			p.used[i] = nil
			delete(p.info, ip)
			return
		}
	}
	//??????????????????
	delete(p.staticTable, ip)
}

//
// ReturnBackById
// @Description:
// @receiver p
// @param id
//
func (p *IPAddressPool) ReturnBackById(id string) {
	for ip := range p.info {
		if p.info[ip].UUID == id {
			for i := range p.used {
				if p.used[i].String() == ip {
					p.used[i] = nil
				}
			}
			delete(p.info, ip)
		}
	}
	for ip := range p.staticTable {
		if p.staticTable[ip].UUID == id {
			for i := range p.used {
				if p.used[i].String() == ip {
					p.used[i] = nil
				}
			}
			delete(p.staticTable, ip)
		}
	}
}
