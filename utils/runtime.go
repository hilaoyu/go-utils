package utils

import (
	"net"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
)

func GetSelfPathFull() string {
	root := os.Args[0]

	// 根据相对路径获取可执行文件的绝对路径
	root, _ = filepath.Abs(root)

	return root
}

func GetSelfPath() string {

	return filepath.Dir(GetSelfPathFull())
}

func RunningOs(osName ...string) string {
	if nil != osName {
		for _, name := range osName {
			sName := name
			if "macos" == sName {
				sName = "darwin"
			}
			if sName == runtime.GOOS {
				return name
			}
		}
		return ""
	}
	return runtime.GOOS
}

func GetSelfIps() map[string]map[string]string {
	ips := map[string]map[string]string{
		"v4": {},
		"v6": {},
	}
	netInterfaces, err := net.Interfaces()

	if nil != err {
		return ips
	}
	v4Ips := map[string]string{}
	v6Ips := map[string]string{}
	for _, netInterface := range netInterfaces {

		netInterfaceIps, err := netInterface.Addrs()

		if nil == err {
			i := 0
			for _, addr := range netInterfaceIps {
				ipNet, isIpNet := addr.(*net.IPNet)

				//是网卡并且不是本地环回网卡
				if isIpNet && !ipNet.IP.IsLoopback() {
					ipv4 := ipNet.IP.To4()
					//能正常转成ipv4
					if ipv4 != nil {
						if _, ok := v4Ips[netInterface.Name]; ok {
							i++
							v4Ips[netInterface.Name+"."+strconv.Itoa(i)] = ipv4.String()
						} else {
							v4Ips[netInterface.Name] = ipv4.String() //addr.String() 0.0.0.0/0
						}
					} else {
						if _, ok := v6Ips[netInterface.Name]; ok {
							i++
							v6Ips[netInterface.Name+"."+strconv.Itoa(i)] = ipNet.IP.String()
						} else {
							v6Ips[netInterface.Name] = ipNet.IP.String() //addr.String() 0.0.0.0/0
						}
					}

				}

			}
		}

	}
	ips["v4"] = v4Ips
	ips["v6"] = v6Ips

	return ips
}

func GetSelfV4Ips() map[string]string {
	ips := GetSelfIps()
	if v4Ips, ok := ips["v4"]; ok {
		return v4Ips
	}

	return nil
}

func GetSelfV4IFaceFirst() string {
	ips := GetSelfV4Ips()
	for iFName, _ := range ips {
		return iFName
	}

	return ""
}
func GetSelfV4IpFirst(iFace ...string) string {
	ips := GetSelfV4Ips()
	iFaceName := ""
	if len(iFace) > 0 {
		iFaceName = iFace[0]
	}
	for iFName, ipAddr := range ips {
		if "" == iFaceName {
			return ipAddr
		} else if iFName == iFaceName {
			return ipAddr
		}
	}

	return ""
}
func GetSelfV6Ips() map[string]string {
	ips := GetSelfIps()
	if v6Ips, ok := ips["v6"]; ok {
		return v6Ips
	}

	return nil
}

func GetSelfV6IFaceFirst() string {
	ips := GetSelfV6Ips()
	for iFName, _ := range ips {
		return iFName
	}

	return ""
}

func GetSelfV6IpFirst(iFace ...string) string {
	ips := GetSelfV6Ips()
	iFaceName := ""
	if len(iFace) > 0 {
		iFaceName = iFace[0]
	}
	for iFName, ipAddr := range ips {
		if "" == iFaceName {
			return ipAddr
		} else if iFName == iFaceName {
			return ipAddr
		}
	}

	return ""
}
