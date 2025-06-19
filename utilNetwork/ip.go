package utilNetwork

import "net"

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
