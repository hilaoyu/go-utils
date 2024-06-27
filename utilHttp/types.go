package utilHttp

import (
	"github.com/hilaoyu/go-utils/utilEnc"
	"github.com/hilaoyu/go-utils/utilLogger"
	"github.com/hilaoyu/go-utils/utilProxy"
	"net/http"
	"net/url"
	"time"
)

type HttpServerListenAddr struct {
	Network string `json:"network,omitempty"`
	Addr    string `json:"addr,omitempty"`
}

type HttpServer struct {
	listenAddresses       []*HttpServerListenAddr
	server                *http.Server
	sslServerCertFile     string
	sslServerKeyFile      string
	sslVerifyClientCaFile string
}

type HttpClient struct {
	timeout                 time.Duration
	baseUrl                 string
	sslVerify               bool
	sslClientCertPemPath    string
	sslClientCertPemContent []byte
	sslClientCertPemKey     []byte
	lastRequestParams       url.Values
	lastRespStatusCode      int
	client                  *http.Client

	rawBody     string
	params      url.Values
	needEncData map[string]interface{}
	headers     map[string]string

	useProxy            string
	proxySocks5         utilProxy.UtilProxy
	proxySocks5Addr     string
	proxySocks5user     string
	proxySocks5Password string

	aesEncryptor *utilEnc.AesEncryptor
	aesEncAppId  string

	logger *utilLogger.Logger
}

type ApiReturnJson struct {
	Status  bool        `json:"status"`
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Errors  []string    `json:"errors"`
	Debug   []string    `json:"debug,omitempty"`
	Data    interface{} `json:"data"`
}
