package utilHttp

import (
	"bytes"
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"github.com/hilaoyu/go-utils/utilLogger"
	"log"
	"net"
	"net/http"
	"os"
	"strings"
	"time"
)

type FilteringLogger struct {
	logger         *utilLogger.Logger
	excludeStrings []string
}

func (fl FilteringLogger) Write(b []byte) (n int, err error) {
	for _, excludeString := range fl.excludeStrings {
		if bytes.Index(b, []byte(excludeString)) > -1 {
			// Filter out the line that matches the pattern
			return len(b), nil
		}
	}
	fl.logger.Error(string(b))
	n = len(b)
	return
}

func NewFilteringWriter(logger *utilLogger.Logger, excludeStrings ...string) *FilteringLogger {
	return &FilteringLogger{
		logger:         logger,
		excludeStrings: excludeStrings,
	}
}

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

func (s *HttpServer) Run(logger *utilLogger.Logger, addresses ...*ServerListenAddr) {
	var err error
	if nil == logger {
		logger = utilLogger.NewLogger()
		_ = logger.AddConsoleWriter()
	}
	s.logger = logger

	s.server.ErrorLog = log.New(
		NewFilteringWriter(
			logger,
			"http: TLS handshake error",
		),
		"http serv: ",
		log.LstdFlags,
	)

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
		certPEMBlock, err1 := os.ReadFile(s.sslVerifyClientCaFile)
		if err1 != nil {
			logger.FatalF("sslVerifyClientCaFile error:%v", err1)
			return
		}
		caPool := x509.NewCertPool()
		caPool.AppendCertsFromPEM(certPEMBlock)
		tlsConfig.ClientCAs = caPool
		s.server.TLSConfig = tlsConfig
	}

	quit := make(chan os.Signal)
	for _, listenAddr := range s.listenAddresses {
		listenAddr.Network = strings.ToLower(listenAddr.Network)
		var listener net.Listener
		listener, err = net.Listen(listenAddr.Network, listenAddr.Addr)
		if nil != err {
			logger.ErrorF("server listen %s://%s , error: %v\n", listenAddr.Network, listenAddr.Addr, err)
			continue
		}
		s.listeners = append(s.listeners, listener)

		if "unix" == listenAddr.Network && listenAddr.Uid > 0 && listenAddr.Gid > 0 {
			if err = os.Chown(listenAddr.Addr, listenAddr.Uid, listenAddr.Gid); err != nil {
				err = fmt.Errorf("server listen %s://%s , Chmod error: %v\n", listenAddr.Network, listenAddr.Addr, err)
				quit <- os.Interrupt
			}
		}
		if "" != listenAddr.SslServerCertFile && "" != listenAddr.SslServerKeyFile {
			logger.InfoF("server serv Tls: %s://%s\n", listenAddr.Network, listenAddr.Addr)
			go func() {
				err = s.server.ServeTLS(listener, listenAddr.SslServerCertFile, listenAddr.SslServerKeyFile)
				if nil != err {
					err = fmt.Errorf("server serv TLS : %s://%s ,errpr: %v", listenAddr.Network, listenAddr.Addr, err)
					quit <- os.Interrupt
				}
			}()
		} else {
			logger.InfoF("server serv : %s://%s\n", listenAddr.Network, listenAddr.Addr)
			go func() {
				err = s.server.Serve(listener)
				err = fmt.Errorf("server serv : %s://%s ,errpr: %v", listenAddr.Network, listenAddr.Addr, err)
				quit <- os.Interrupt
			}()

		}

		if nil != err && err != http.ErrServerClosed {
			logger.ErrorF("%v\n", err)
		}

	}

	if nil != err {
		logger.ErrorF("%v\n", err)
	}

	<-quit
	s.Shutdown()
	return
}
func (s *HttpServer) Shutdown() {
	logger := s.logger
	if nil == logger {
		logger = utilLogger.NewLogger()
		_ = logger.AddConsoleWriter()
	}

	logger.Info("Shutdown Server ...")
	var err error
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(10)*time.Second)
	defer func() {
		for _, listener := range s.listeners {
			if nil != listener {
				logger.InfoF("close  listener %s", listener.Addr())
				_ = listener.Close()

			}
		}
		cancel()
	}()
	if err = s.server.Shutdown(ctx); err != nil {
		logger.FatalF("Server Shutdown: %v", err)
	}
	logger.Info("Server exiting")
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
