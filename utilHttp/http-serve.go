package utilHttp

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"github.com/hilaoyu/go-utils/utilLogger"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"time"
)

func NewHttpServe(handler http.Handler, addresses ...*ServerListenAddr) (s *HttpServer) {
	s = &HttpServer{server: &http.Server{Handler: handler}, listenAddresses: addresses}
	return
}

func (s *HttpServer) SetReadTimeout(t time.Duration) *HttpServer {
	s.server.ReadTimeout = t
	return s
}
func (s *HttpServer) SetReadHeaderTimeout(t time.Duration) *HttpServer {
	s.server.ReadHeaderTimeout = t
	return s
}
func (s *HttpServer) SetWriteTimeout(t time.Duration) *HttpServer {
	s.server.WriteTimeout = t
	return s
}
func (s *HttpServer) SetIdleTimeout(t time.Duration) *HttpServer {
	s.server.IdleTimeout = t
	return s
}

func (s *HttpServer) SetMaxHeaderBytes(i int) *HttpServer {
	s.server.MaxHeaderBytes = i
	return s
}

func (s *HttpServer) VerifyClientSsl(caFile string) *HttpServer {
	s.sslVerifyClientCaFile = caFile
	return s
}

func (s *HttpServer) Run(logger *utilLogger.Logger, addresses ...*ServerListenAddr) (err error) {
	if nil == logger {
		logger = utilLogger.NewLogger()
		logger.AddConsoleWriter()
	}

	if len(addresses) > 0 {
		s.listenAddresses = append(s.listenAddresses, addresses...)
	}

	if len(s.listenAddresses) <= 0 {
		logger.Fatal("listen addresses is empty")
		return
	}

	if "" != s.sslVerifyClientCaFile {
		tlsConfig := &tls.Config{}
		tlsConfig.ClientAuth = tls.RequireAndVerifyClientCert
		certPEMBlock, err := os.ReadFile(s.sslVerifyClientCaFile)
		if err != nil {
			logger.FatalF("sslVerifyClientCaFile error:%v", err)
			return err
		}
		caPool := x509.NewCertPool()
		caPool.AppendCertsFromPEM(certPEMBlock)
		tlsConfig.ClientCAs = caPool
		s.server.TLSConfig = tlsConfig
	}

	for _, listenAddr := range s.listenAddresses {
		listenAddr.Network = strings.ToLower(listenAddr.Network)
		listener, err := net.Listen(listenAddr.Network, listenAddr.Addr)
		if nil != err {
			logger.ErrorF("server listen %s://%s , error: %v\n", listenAddr.Network, listenAddr.Addr, err)
			continue
		}
		if "unix" == listenAddr.Network && listenAddr.Uid > 0 && listenAddr.Gid > 0 {
			if err = os.Chown(listenAddr.Addr, listenAddr.Uid, listenAddr.Gid); err != nil {
				logger.ErrorF("server listen %s://%s , Chmod error: %v\n", listenAddr.Network, listenAddr.Addr, err)
				listener.Close()
				continue
			}
		}
		if "" != listenAddr.SslServerCertFile && "" != listenAddr.SslServerKeyFile {
			go func() {
				err = s.server.ServeTLS(listener, listenAddr.SslServerCertFile, listenAddr.SslServerKeyFile)
			}()
		} else {
			go func() {
				err = s.server.Serve(listener)
			}()

		}

		if nil != err && err != http.ErrServerClosed {
			logger.ErrorF("server serv %s://%s , error: %v\n", listenAddr.Network, listenAddr.Addr, err)
		} else {
			logger.InfoF("server listen: %s://%s\n", listenAddr.Network, listenAddr.Addr)
		}

	}

	// 等待中断信号以优雅地关闭服务器（设置 5 秒的超时时间）
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
	logger.Info("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(10)*time.Second)
	defer cancel()
	if err = s.server.Shutdown(ctx); err != nil {
		logger.FatalF("Server Shutdown: %v", err)
	}
	logger.Info("Server exiting")
	return
}

func GetClientIps(r *http.Request) (ips []string) {

	ip := strings.TrimSpace(strings.Split(r.Header.Get("X-Forwarded-For"), ",")[0])
	if ip != "" {
		ips = append(ips, ip)
	}

	ip = strings.TrimSpace(r.Header.Get("X-Real-Ip"))
	if ip != "" {
		ips = append(ips, ip)
	}

	if ip, _, err := net.SplitHostPort(strings.TrimSpace(r.RemoteAddr)); err == nil {
		ips = append(ips, ip)
	}

	return
}
func GetClientIP(r *http.Request) (ip string) {
	ips := GetClientIps(r)
	if len(ips) > 0 {
		ip = ips[0]
	}
	return
}
func GetClientPublicIP(r *http.Request) (ip string) {
	ips := GetClientIps(r)
	for _, ipTemp := range ips {
		ipParse := net.ParseIP(ipTemp)
		if nil != ipParse && ipParse.IsGlobalUnicast() && !ipParse.IsPrivate() {
			ip = ipTemp
			return
		}
	}
	return
}
func GetClientPlatform(r *http.Request) (clientPlatform string) {
	clientPlatform = r.Header.Get("ClientPlatform")
	if "" == clientPlatform {
		clientPlatform = "web"
	}
	return
}
func GetClientSource(r *http.Request) (clientSource string) {
	clientSource = r.Header.Get("ClientSource")
	if "" == clientSource {
		clientSource = GetClientIP(r)
	}
	return
}
