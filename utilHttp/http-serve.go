package utilHttp

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

type HttpServer struct {
	server                *http.Server
	sslServerCertFile     string
	sslServerKeyFile      string
	sslVerifyClientCaFile string
}

func NewHttpServe(addr string, handler http.Handler) (s *HttpServer) {
	s = &HttpServer{server: &http.Server{Addr: addr, Handler: handler}}
	return
}

func (s *HttpServer) SetReadTimeout(t time.Duration) {
	s.server.ReadTimeout = t
}
func (s *HttpServer) SetReadHeaderTimeout(t time.Duration) {
	s.server.ReadHeaderTimeout = t
}
func (s *HttpServer) SetWriteTimeout(t time.Duration) {
	s.server.WriteTimeout = t
}
func (s *HttpServer) SetIdleTimeout(t time.Duration) {
	s.server.IdleTimeout = t
}

func (s *HttpServer) SetMaxHeaderBytes(i int) {
	s.server.MaxHeaderBytes = i
}

func (s *HttpServer) UseServerSsl(certFile string, keyFile string) {
	s.sslServerCertFile = certFile
	s.sslServerKeyFile = keyFile
}

func (s *HttpServer) VerifyClientSsl(caFile string) {
	s.sslVerifyClientCaFile = caFile
}

func (s *HttpServer) Run(logger *log.Logger) (err error) {

	if "" != s.sslServerCertFile || "" != s.sslServerKeyFile {
		tlsConfig := &tls.Config{}

		tlsConfig.Certificates = make([]tls.Certificate, 1)
		tlsConfig.Certificates[0], err = tls.LoadX509KeyPair(s.sslServerCertFile, s.sslServerKeyFile)
		if err != nil {
			return err
		}

		if "" != s.sslVerifyClientCaFile {
			tlsConfig.ClientAuth = tls.RequireAndVerifyClientCert
			certPEMBlock, err := os.ReadFile(s.sslVerifyClientCaFile)
			if err != nil {
				return err
			}
			caPool := x509.NewCertPool()
			caPool.AppendCertsFromPEM(certPEMBlock)
			tlsConfig.ClientCAs = caPool
		}

		s.server.TLSConfig = tlsConfig
	}

	go func() {
		// 服务连接
		if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatalf("listen: %s\n", err)
		}
	}()

	// 等待中断信号以优雅地关闭服务器（设置 5 秒的超时时间）
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
	logger.Println("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err = s.server.Shutdown(ctx); err != nil {
		logger.Fatal("Server Shutdown:", err)
	}
	logger.Println("Server exiting")
	return
}
