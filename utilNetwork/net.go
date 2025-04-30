package utilNetwork

import (
	"github.com/c-robinson/iplib"
	"math"
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
func (un *UtilNet) IpCount() (count uint32) {
	size, bl := un.Mask().Size()
	count = uint32(math.Pow(2, float64(bl-size)))
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

func (un *UtilNet) NetworkAddress() net.IP {
	if net4, ok := un.Net.(iplib.Net4); ok {
		return net4.NetworkAddress()
	}
	return nil
}
func (un *UtilNet) BroadcastAddress() net.IP {
	if net4, ok := un.Net.(iplib.Net4); ok {
		return net4.BroadcastAddress()
	}
	return nil
}

func (un *UtilNet) AvailableIps() (ips []net.IP) {
	lastIp := un.LastAddress()
	currentIp := iplib.NextIP(un.IP())
	ips = append(ips, currentIp)
	for {
		currentIp = iplib.NextIP(currentIp)
		ips = append(ips, currentIp)
		if currentIp.Equal(lastIp) {
			break
		}
	}
	return
}
