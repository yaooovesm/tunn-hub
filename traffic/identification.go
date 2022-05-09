package traffic

import (
	"github.com/songgao/water/waterutil"
	"strconv"
	"strings"
)

type Way int

const (
	In      Way = 0
	Out     Way = 1
	Default Way = 2
)

//
// Identification
// @Description:
// @param b
// @param way
// @return string
//
func Identification(b []byte, way Way) string {
	defer func() {
		if err := recover(); err != nil {
			return
		}
	}()
	src := ""
	dst := ""
	if waterutil.IPv4Protocol(b) == waterutil.TCP || waterutil.IPv4Protocol(b) == waterutil.UDP {
		srcIp := waterutil.IPv4Source(b)
		dstIp := waterutil.IPv4Destination(b)
		srcPort := waterutil.IPv4SourcePort(b)
		dstPort := waterutil.IPv4DestinationPort(b)
		src = strings.Join([]string{srcIp.To4().String(), strconv.FormatUint(uint64(srcPort), 10)}, ":")
		dst = strings.Join([]string{dstIp.To4().String(), strconv.FormatUint(uint64(dstPort), 10)}, ":")
	} else {
		src = waterutil.IPv4Source(b).String()
		dst = waterutil.IPv4Destination(b).String()
	}
	if src == "" || dst == "" {
		return ""
	}
	switch way {
	case Out:
		return dst + "->" + src
	case In, Default:
		return src + "->" + dst
	}
	return ""
}

//
// IdentificationV2
// @Description: In Default -> |protocol|src,src_port(tcp or udp)|dst,dst_port(tcp or udp)|
//				 Out        -> |protocol|dst,dst_port(tcp or udp)|src,src_port(tcp or udp)|
// @param b
// @param way
// @return string
//
func IdentificationV2(b []byte, way Way) string {
	defer func() {
		if err := recover(); err != nil {
			return
		}
	}()
	if len(b) < 20 {
		return ""
	}
	//src := []byte{b[9], b[12], b[13], b[14], b[15], byte(0)}
	//dst := []byte{b[9], b[16], b[17], b[18], b[19], byte(0)}
	//remove protocol
	src := []byte{byte(0), b[12], b[13], b[14], b[15], byte(0)}
	dst := []byte{byte(0), b[16], b[17], b[18], b[19], byte(0)}
	//9->protocol
	//12:16->src
	//16:20->dst
	switch b[9] {
	case 0x11, 0x06:
		//tcp,udp
		//ihl := b[0] & 0x0F
		//pl := b[ihl*4:]
		//src = append(src, pl[0:2]...)
		//dst = append(dst, pl[2:4]...)
		src[5] = byte(waterutil.IPv4SourcePort(b))
		dst[5] = byte(waterutil.IPv4DestinationPort(b))
	}
	switch way {
	case Out:
		return string(dst) + string(src)
	case In, Default:
		return string(src) + string(dst)
	}
	return ""
}
