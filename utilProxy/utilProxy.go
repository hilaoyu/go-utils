package utilProxy

import (
	"context"
	"golang.org/x/net/proxy"
	"net"
	"strings"
	"time"
)

type UtilProxy interface {
	NewConn(network string, addr string, timeout ...time.Duration) (conn net.Conn, err error)
	Dial(network string, addr string) (conn net.Conn, err error)
	DialContext(ctx context.Context, network string, addr string) (conn net.Conn, err error)
}

type ProxySocks5 struct {
	dialer proxy.Dialer
}

func NewProxySocks5(proxyAddress string, proxyUser string, proxyPassword string) (utilProxy UtilProxy, err error) {
	var proxyAuth *proxy.Auth
	proxyUser = strings.TrimSpace(proxyUser)
	if "" != proxyUser {
		proxyAuth = &proxy.Auth{
			User:     proxyUser,
			Password: proxyPassword,
		}
	}
	proxyDialer, err := proxy.SOCKS5("tcp", proxyAddress, proxyAuth, proxy.Direct)
	if nil != err {
		return
	}

	utilProxy = &ProxySocks5{dialer: proxyDialer}

	return
}

func (p *ProxySocks5) Dial(network string, addr string) (conn net.Conn, err error) {
	return p.dialer.Dial(network, addr)
}
func (p *ProxySocks5) DialContext(ctx context.Context, network string, addr string) (conn net.Conn, err error) {
	if dialer, ok := p.dialer.(proxy.ContextDialer); ok {
		return dialer.DialContext(ctx, network, addr)
	}
	return p.dialer.Dial(network, addr)
}
func (p *ProxySocks5) NewConn(network string, addr string, timeout ...time.Duration) (conn net.Conn, err error) {
	connTimeout := time.Duration(10) * time.Second
	if len(timeout) > 0 && timeout[0] > 0 {
		connTimeout = timeout[0]
	}
	dialer, ok := p.dialer.(proxy.ContextDialer)
	if connTimeout > 0 && ok {
		ctx, cancel := context.WithTimeout(context.Background(), connTimeout)
		defer cancel()
		conn, err = dialer.DialContext(ctx, network, addr)
	} else {
		conn, err = p.dialer.Dial(network, addr)
	}

	return
}
