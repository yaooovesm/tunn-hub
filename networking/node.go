package networking

import "net"

//
// TableNode
// @Description:
//
type TableNode struct {
	Name       string
	Net        *net.IPNet
	TunnelAddr string
	Next       *TableNode
	Last       *TableNode
	rank       int
	symbol     net.IP
}

//
// NewTableNode
// @Description:
// @param dst
// @param tunnel
// @return node
// @return err
//
func NewTableNode(dst string, uuid string) (node *TableNode, err error) {
	ip, ipNet, err := net.ParseCIDR(dst)
	if err != nil {
		return nil, err
	}
	name := ipNet.String()
	return &TableNode{
		Name:       name,
		Net:        ipNet,
		TunnelAddr: uuid,
		symbol:     ip.To4(),
	}, nil
}

//
// NewTableNodeWithLimit
// @Description:
// @param dst
// @param tunnel
// @return node
// @return err
//
func NewTableNodeWithLimit(dst string, uuid string, minimumMaskSize int) (node *TableNode, err error) {
	ip, ipNet, err := net.ParseCIDR(dst)
	if err != nil {
		return nil, err
	}
	if size, _ := ipNet.Mask.Size(); size < minimumMaskSize {
		return nil, ErrMaskSizeOutOfLimit
	}
	name := ipNet.String()
	return &TableNode{
		Name:       name,
		Net:        ipNet,
		TunnelAddr: uuid,
		symbol:     ip.To4(),
	}, nil
}
