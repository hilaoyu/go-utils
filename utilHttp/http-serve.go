package utilHttp

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"github.com/hilaoyu/go-utils/utilLogger"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"time"
)

func NewHttpServe(handler http.Handler, addresses ...*HttpServerListenAddr) (s *HttpServer) {
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

func (s *HttpServer) UseServerSsl(certFile string, keyFile string) *HttpServer {
	s.sslServerCertFile = certFile
	s.sslServerKeyFile = keyFile
	return s
}

func (s *HttpServer) VerifyClientSsl(caFile string) *HttpServer {
	s.sslVerifyClientCaFile = caFile
	return s
}

func (s *HttpServer) Run(logger *utilLogger.Logger) (err error) {

	if "" != s.sslVerifyClientCaFile {
		tlsConfig := &tls.Config{}
		tlsConfig.ClientAuth = tls.RequireAndVerifyClientCert
		certPEMBlock, err := os.ReadFile(s.sslVerifyClientCaFile)
		if err != nil {
			return err
		}
		caPool := x509.NewCertPool()
		caPool.AppendCertsFromPEM(certPEMBlock)
		tlsConfig.ClientCAs = caPool
		s.server.TLSConfig = tlsConfig
	}

	servErrorChan := make(chan error, 100)
	for _, listenAddr := range s.listenAddresses {
		listenAddr.Network = strings.ToLower(listenAddr.Network)
		httpListenScheme := listenAddr.Network
		listener, err := net.Listen(listenAddr.Network, listenAddr.Addr)
		if nil != err {
			logger.Fatal(fmt.Sprintf("http server listen %s://%s , error: %v\n", httpListenScheme, listenAddr.Addr, err))
			continue
		}
		switch listenAddr.Network {
		case "tcp", "udp":
			// 服务连接
			if "" != s.sslServerCertFile || "" != s.sslServerKeyFile {
				httpListenScheme = "https"
				go func() {
					err1 := s.server.ServeTLS(listener, s.sslServerCertFile, s.sslServerKeyFile)
					servErrorChan <- fmt.Errorf("http server serv %s://%s , error: %v\n", httpListenScheme, listenAddr.Addr, err1)
				}()
			} else {
				httpListenScheme = "http"
				go func() {
					err1 := s.server.Serve(listener)
					servErrorChan <- fmt.Errorf("http server serv %s://%s , error: %v\n", httpListenScheme, listenAddr.Addr, err1)
				}()

			}
			break
		default:
			go func() {
				err1 := s.server.Serve(listener)
				servErrorChan <- fmt.Errorf("http server serv %s://%s , error: %v\n", httpListenScheme, listenAddr.Addr, err1)
			}()
			break

		}

		if err != nil && err != http.ErrServerClosed {
			logger.Fatal(fmt.Sprintf("http server listen %s://%s , error: %v\n", httpListenScheme, listenAddr.Addr, err))
		} else {
			logger.Info(fmt.Sprintf("http server listen: %s://%s\n", httpListenScheme, listenAddr.Addr))
		}

	}

	go func() {
		errServ := <-servErrorChan
		logger.Fatal(fmt.Sprintf(" %v\n", errServ))
	}()
	// 等待中断信号以优雅地关闭服务器（设置 5 秒的超时时间）
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
	logger.Info("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(10)*time.Second)
	defer cancel()
	if err = s.server.Shutdown(ctx); err != nil {
		logger.Fatal(fmt.Sprintf("Server Shutdown:", err))
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
