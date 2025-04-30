package utilNetwork

import (
	"fmt"
	"github.com/bits-and-blooms/bitset"
	"net"
	"sync"
)

type UtilIpAm struct {
	net    *UtilNet
	ipSet  *bitset.BitSet
	locker sync.Mutex
}

func NewUtilIpAm(netCidr string) (ipAm *UtilIpAm, err error) {
	utilNet, err := NewUtilNet(netCidr)
	if nil != err {
		return
	}
	netIpCount := uint(utilNet.IpCount())
	ipSet := bitset.New(netIpCount)
	if nil != utilNet.NetworkAddress() {
		ipSet.Set(0) //网络地址设置为已使用
	}

	if nil != utilNet.BroadcastAddress() {
		ipSet.Set(netIpCount - 1) //广播地址设置为已使用
	}

	ipAm = &UtilIpAm{
		net:   utilNet,
		ipSet: ipSet,
	}

	return
}

func (ia *UtilIpAm) FindAvailableIp() (ip net.IP, err error) {
	pos := uint(0)
	found := false
	pos, found = ia.ipSet.NextClear(pos)
	if !found {
		err = fmt.Errorf("no ip available")
		return
	}
	ip = net.ParseIP(ia.net.GetIpByPosition(uint32(pos)))
	return
}
func (ia *UtilIpAm) FindAvailableIpAndUse() (ip net.IP, err error) {
	ia.locker.Lock()
	defer func() {
		ia.locker.Unlock()
	}()
	ip, err = ia.FindAvailableIp()
	if nil != err {
		return
	}
	ia.UseIp(ip)
	return
}
func (ia *UtilIpAm) UseIp(ip net.IP) {
	if nil == ip {
		return
	}
	pos := ia.net.IpPosition(ip)
	ia.ipSet.Set(uint(pos))
}
func (ia *UtilIpAm) UseIpStr(ipStr string) {
	ia.UseIp(net.ParseIP(ipStr))
}
func (ia *UtilIpAm) UnUseIp(ip net.IP) {
	if nil == ip {
		return
	}
	pos := ia.net.IpPosition(ip)
	ia.ipSet.Clear(uint(pos))
}
func (ia *UtilIpAm) UnUseIpStr(ipStr string) {
	ia.UnUseIp(net.ParseIP(ipStr))
}

func (ia *UtilIpAm) UsedIpCount() uint32 {
	return uint32(ia.ipSet.Count())
}
func (ia *UtilIpAm) AvailableIpCount() uint32 {
	return uint32(ia.ipSet.Len() - ia.ipSet.Count())
}
