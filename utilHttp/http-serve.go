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

func (s *HttpServer) Run(logger *log.Logger) (err error) {

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

	go func() {
		// 服务连接
		if "" != s.sslServerCertFile || "" != s.sslServerKeyFile {
			err = s.server.ListenAndServeTLS(s.sslServerCertFile, s.sslServerKeyFile)
		} else {
			err = s.server.ListenAndServe()
		}

		if err != nil && err != http.ErrServerClosed {
			logger.Fatalf("listen: %s\n", err)
		}
	}()

	// 等待中断信号以优雅地关闭服务器（设置 5 秒的超时时间）
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
	logger.Println("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(10)*time.Second)
	defer cancel()
	if err = s.server.Shutdown(ctx); err != nil {
		logger.Fatal("Server Shutdown:", err)
	}
	logger.Println("Server exiting")
	return
}
