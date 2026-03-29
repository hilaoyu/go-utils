package utilNetwork

import "net"

func Parse(addr string)(ip net.IP, ipNet *net.IPNet)  {
	ip,ipNet,err := net.ParseCIDR(addr)
	if nil != err {
		ip = net.ParseIP(addr)
	}
	return
}

func IsPublicIP(ip net.IP) bool {
	if ip == nil {
		return false
	}

	if ip.IsPrivate() || ip.IsLoopback() || ip.IsLinkLocalMulticast() || ip.IsLinkLocalUnicast() {
		return false
	}
	return true
}

func IpInCidr(ip string, cidr string) bool {
	if "" == ip || "" == cidr {
		return false
	}
	ipParse := net.ParseIP(ip)
	_, netParse, err := net.ParseCIDR(cidr)
	if nil != err {
		return false
	}

	if netParse.Contains(ipParse) {
		return true
	}

	return false
}

func IpV4ToInt(ipStr string) uint32 {
	ip,_ := Parse(ipStr)
	if ip == nil {
		return 0
	}
	ip = ip.To4()
	if ip == nil {
		return 0
	}

	return uint32(ip[0])<<24 |
		uint32(ip[1])<<16 |
		uint32(ip[2])<<8 |
		uint32(ip[3])

}


func IntToIPV4(n uint32) string {
	ip := make(net.IP, 4)
	ip[0] = byte(n >> 24)
	ip[1] = byte(n >> 16)
	ip[2] = byte(n >> 8)
	ip[3] = byte(n)
	return ip.String()
}