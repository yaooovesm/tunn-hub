package utils

import (
	"net"
	"strconv"
	"strings"
)

//
// CIDR2Mask
// @Description:
// @param mask
// @return string
//
func CIDR2Mask(mask net.IPMask) string {
	val := make([]byte, len(mask))
	copy(val, mask)

	var s []string
	for _, i := range val[:] {
		s = append(s, strconv.Itoa(int(i)))
	}
	return strings.Join(s, ".")
}
