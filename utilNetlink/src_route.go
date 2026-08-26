package utilNetlink

import (
	"fmt"
	"net"
	"strings"

	"github.com/hilaoyu/go-utils/utilNetwork"
	"github.com/vishvananda/netlink"
)

func ErrIsNoSuch(err error) bool {
	return strings.Contains(err.Error(), "no such")
}

func SrcRouteAdd(src string, gateway string) (err error) {
	gateway = strings.TrimSpace(gateway)
	src = strings.TrimSpace(src)
	if "" == gateway || "" == src {
		return
	}

	table := int(utilNetwork.IpV4ToInt(src))
	//fmt.Println(table)

	if table <= 0 {
		err = fmt.Errorf("table 错误")
		return
	}

	route, err := srcRouteBuildNetlinkRoute(gateway, table)
	if nil != err {
		return
	}

	_ = netlink.RouteDel(route)

	defer func() {
		if nil != err {
			_ = netlink.RouteDel(route)
		}
	}()
	err = netlink.RouteAdd(route)
	if nil != err {
		return
	}

	rule, err := srcRouteBuildNetlinkRule(src, table)
	if nil != err {
		return
	}

	_ = netlink.RuleDel(rule)
	defer func() {
		if nil != err {
			_ = netlink.RuleDel(rule)
		}
	}()

	err = netlink.RuleAdd(rule)
	if nil != err {
		return
	}

	return
}
func SrcRouteDelete(src string) (err error) {
	src = strings.TrimSpace(src)
	if "" == src {
		return
	}

	table := int(utilNetwork.IpV4ToInt(src))
	if table <= 0 {
		err = fmt.Errorf("table 错误")
		return
	}
	route, err := srcRouteBuildNetlinkRoute("", table)
	if nil != err {
		return
	}

	err = netlink.RouteDel(route)
	if nil != err {
		if ErrIsNoSuch(err) {
			err = nil
		} else {
			return
		}
	}

	rule, err := srcRouteBuildNetlinkRule(src, table)
	if nil != err {
		return
	}

	err = netlink.RuleDel(rule)

	if nil != err {
		if ErrIsNoSuch(err) {
			err = nil
		} else {
			return
		}
	}

	return
}

func srcRouteBuildNetlinkRoute(gateway string, table int) (route *netlink.Route, err error) {
	gateway = strings.TrimSpace(gateway)

	if table <= 0 {
		err = fmt.Errorf("table 错误")
		return
	}

	_, detNet := utilNetwork.Parse("0.0.0.0/0")

	route = &netlink.Route{
		Dst:   detNet,
		Table: table,
		/*Via: &netlink.Via{
			AddrFamily: 0,
			Addr:       gatewayIp,
		},*/
	}
	if "" != gateway {
		gatewayIp, _ := utilNetwork.Parse(gateway)
		if nil == gatewayIp {
			err = fmt.Errorf("gateway 错误")
			return
		}
		route.Gw = gatewayIp
	}
	return
}
func srcRouteBuildNetlinkRule(src string, table int) (rule *netlink.Rule, err error) {
	src = strings.TrimSpace(src)
	if "" == src {
		err = fmt.Errorf("src 为空")
		return
	}

	if table <= 0 {
		err = fmt.Errorf("table 错误")
		return
	}
	srcIp, srcNet := utilNetwork.Parse(src)

	if nil == srcNet {
		if nil == srcIp {
			err = fmt.Errorf("src 错误")
			return
		}
		srcNet = &net.IPNet{
			IP:   srcIp,
			Mask: net.CIDRMask(32, 32), // 单一 IP = /32
		}
	}
	//_, detNet := utilNetwork.Parse("0.0.0.0/0")
	rule = netlink.NewRule()
	rule.Src = srcNet
	rule.Table = table
	return
}
