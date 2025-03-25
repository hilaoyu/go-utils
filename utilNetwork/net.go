package utilNetwork

import (
	"github.com/c-robinson/iplib"
	"net"
)

type UtilNet struct {
	iplib.Net
}

func NewUtilNet(cidr string) (utilNet *UtilNet, err error) {
	_, vnet, err := iplib.ParseCIDR(cidr)
	if nil != err {
		return
	}

	utilNet = &UtilNet{Net: vnet}
	return
}

func (un *UtilNet) IpPosition(ip net.IP) uint32 {
	return iplib.DeltaIP(un.IP(), ip)
}
func (un *UtilNet) MaskSize() (size int) {
	size, _ = un.Mask().Size()
	return
}
func (un *UtilNet) GetIpByPosition(position ...uint32) string {
	var p = uint32(0)
	if len(position) > 0 {
		p = position[0]
	}
	if p > 0 {
		p = p - 1
	}
	return iplib.IncrementIPBy(un.FirstAddress(), p).String()
}
func (un *UtilNet) GetIpByPositionReverse(position ...uint32) string {
	var p = uint32(0)
	if len(position) > 0 {
		p = position[0]
	}
	if p > 0 {
		p = p - 1
	}
	return iplib.DecrementIPBy(un.LastAddress(), p).String()
}
func (un *UtilNet) GetNetMask() string {
	return iplib.HexStringToIP(un.Mask().String()).String()
}
