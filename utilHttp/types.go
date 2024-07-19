package utilHttp

import (
	"github.com/hilaoyu/go-utils/utilEnc"
	"github.com/hilaoyu/go-utils/utilLogger"
	"github.com/hilaoyu/go-utils/utilProxy"
	"net/http"
	"net/url"
	"time"
)

type ServerListenAddr struct {
	Network           string `json:"network,omitempty"`
	Addr              string `json:"addr,omitempty"`
	Uid               int    `json:"uid,omitempty"`
	Gid               int    `json:"gid,omitempty"`
	SslServerCertFile string `json:"ssl_server_cert_file,omitempty"`
	SslServerKeyFile  string `json:"ssl_server_key_file,omitempty"`
}

type HttpServer struct {
	listenAddresses       []*ServerListenAddr
	server                *http.Server
	sslVerifyClientCaFile string
}

type HttpClient struct {
	timeout                 time.Duration
	baseUrl                 string
	sslVerify               bool
	sslClientCertPemPath    string
	sslClientCertPemContent []byte
	sslClientCertPemKey     []byte
	lastRequestUrl          string
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

type ApiDataJson struct {
	Status  bool        `json:"status"`
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Errors  []string    `json:"errors"`
	Debug   []string    `json:"debug,omitempty"`
	Data    interface{} `json:"data"`
}

type ApiDataSelectOption struct {
	Key      interface{} `json:"key,omitempty"`
	Label    string      `json:"label,omitempty"`
	Value    interface{} `json:"value,omitempty"`
	Selected bool        `json:"selected,omitempty"`
	Disabled bool        `json:"disabled,omitempty"`
}
