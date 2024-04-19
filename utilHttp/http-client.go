package utilHttp

import (
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/hilaoyu/go-utils/utilEnc"
	"github.com/hilaoyu/go-utils/utilProxy"
	"github.com/hilaoyu/go-utils/utilSsl"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

const PROXY_TYPE_SOCKS5 = "socks5"

func NewHttpClient(baseUrl string, timeout ...time.Duration) (uh *HttpClient) {
	uh = &HttpClient{
		timeout:                 0,
		baseUrl:                 baseUrl,
		headers:                 map[string]string{},
		sslVerify:               true,
		sslClientCertPemPath:    "",
		sslClientCertPemContent: nil,
		sslClientCertPemKey:     nil,
		client:                  nil,
	}
	if len(timeout) > 0 {
		uh.timeout = timeout[0]
	}

	return uh.buildClient()
}

func (uh *HttpClient) UseProxySocks5(proxyAddr string, proxyUser string, proxyPassword string) *HttpClient {
	proxyAddr = strings.TrimSpace(proxyAddr)
	uh.useProxy = PROXY_TYPE_SOCKS5
	if "" != proxyAddr {
		uh.proxySocks5Addr = proxyAddr
		uh.proxySocks5user = proxyUser
		uh.proxySocks5Password = proxyPassword

		d, err := utilProxy.NewProxySocks5(uh.proxySocks5Addr, uh.proxySocks5user, uh.proxySocks5Password)
		if nil == err {
			uh.proxySocks5 = d
		}

	}

	return uh.buildClient()
}

func (uh *HttpClient) WithAesEncryptor(secret string, appId string) *HttpClient {
	uh.aesEncryptor = utilEnc.NewAesEncryptor(secret)
	uh.aesEncAppId = appId
	return uh.buildClient()
}
func (uh *HttpClient) GetAesEncryptor() *utilEnc.AesEncryptor {
	return uh.aesEncryptor
}

func (uh *HttpClient) buildClient() *HttpClient {
	uh.client = nil

	tlsConfig := &tls.Config{InsecureSkipVerify: !uh.sslVerify}

	if nil == uh.sslClientCertPemContent && "" != uh.sslClientCertPemPath {
		uh.sslClientCertPemKey, uh.sslClientCertPemContent = utilSsl.ParsePemCertFile(uh.sslClientCertPemPath)
	}
	if nil != uh.sslClientCertPemContent {
		tlsCertificate, err := tls.X509KeyPair(uh.sslClientCertPemContent, uh.sslClientCertPemKey)
		if nil == err {
			tlsConfig.Certificates = []tls.Certificate{tlsCertificate}
		} else {
			fmt.Println("tls.X509KeyPair err:", err)
		}
	}

	tr := &http.Transport{
		TLSClientConfig: tlsConfig,
	}

	switch uh.useProxy {
	case PROXY_TYPE_SOCKS5:
		if nil != uh.proxySocks5 {
			tr.DialContext = uh.proxySocks5.DialContext
		}
		break
	default:
		break
	}

	uh.client = &http.Client{
		Timeout:   uh.timeout,
		Transport: tr,
	}
	return uh
}

func (uh *HttpClient) SetBaseUrl(baseUrl string) *HttpClient {
	uh.baseUrl = baseUrl
	return uh.buildClient()
}
func (uh *HttpClient) SetTimeout(timeout time.Duration) *HttpClient {
	uh.timeout = timeout
	return uh.buildClient()
}
func (uh *HttpClient) SetClientCertPemPath(path string) *HttpClient {
	uh.sslClientCertPemPath = path
	return uh.buildClient()
}

func (uh *HttpClient) SetClientCertPemContent(pemContent []byte, key []byte) *HttpClient {
	uh.sslClientCertPemContent = pemContent
	uh.sslClientCertPemKey = key
	return uh.buildClient()
}

func (uh *HttpClient) SetSslVerify(v bool) *HttpClient {
	uh.sslVerify = v
	return uh.buildClient()
}

func (uh *HttpClient) BasicAuth(user string, password string) *HttpClient {
	uh.AddHeader("Authorization", fmt.Sprintf("Basic %s", base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", user, password)))))
	return uh.buildClient()
}

func (uh *HttpClient) BuildRemoteUrlAndParams(method string, path string) (remoteUrl string, params url.Values, err error) {
	method = strings.ToUpper(method)
	remoteUrl = path
	if "" != uh.baseUrl {
		remoteUrl = strings.TrimRight(uh.baseUrl, "/") + "/" + strings.TrimLeft(path, "/")
	}

	params = uh.params
	if nil != uh.aesEncryptor {

		needEncData := uh.needEncData
		if len(needEncData) > 0 {
			needEncData["_timestamp"] = time.Now().UTC().Unix()
			needEncData["_data_id"] = strconv.FormatInt(time.Now().UTC().UnixNano(), 10)
			enData, err1 := uh.aesEncryptor.Encrypt(needEncData)
			if nil != err1 {
				err = err1
				return
			}
			params.Set("data", enData)
		}
		if "" != uh.aesEncAppId {
			params.Set("app_id", uh.aesEncAppId)
		}
	}

	if "GET" == method || "DELETE" == method {
		urlParse, err1 := url.Parse(remoteUrl)
		if nil != err1 {
			err = err1
			return
		}
		query := urlParse.Query()
		for qk, _ := range params {
			query.Set(qk, params.Get(qk))
		}
		urlParse.RawQuery = query.Encode()

		remoteUrl = urlParse.String()
		params = url.Values{}
	}
	return
}

func (uh *HttpClient) AddHeader(k string, v string) *HttpClient {
	k = strings.TrimSpace(k)
	v = strings.TrimSpace(v)
	if "" != k {
		uh.headers[k] = v
	}
	return uh.buildClient()
}
func (uh *HttpClient) ClearHeader() *HttpClient {
	uh.headers = map[string]string{}
	return uh.buildClient()
}

func (uh *HttpClient) WithParams(params url.Values) *HttpClient {
	uh.params = params
	return uh.buildClient()
}
func (uh *HttpClient) ClearParams() *HttpClient {
	uh.params = url.Values{}
	return uh.buildClient()
}

func (uh *HttpClient) WithEncryptData(data map[string]interface{}) *HttpClient {
	uh.needEncData = data
	return uh.buildClient()
}
func (uh *HttpClient) ClearEncryptData() *HttpClient {
	uh.needEncData = map[string]interface{}{}
	return uh.buildClient()
}

func (uh *HttpClient) GetLastRespStatusCode() int {
	return uh.lastRespStatusCode
}
func (uh *HttpClient) GetLastRequestParams() url.Values {
	return uh.lastRequestParams
}

func (uh *HttpClient) Request(method string, path string, additionalHeaders map[string]string) (resp *http.Response, err error) {
	method = strings.ToUpper(method)

	remoteUrl, params, err := uh.BuildRemoteUrlAndParams(method, path)

	//fmt.Println(remoteUrl, params.Encode())

	req, err := http.NewRequest(method, remoteUrl, strings.NewReader(params.Encode()))

	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	for hk, hv := range uh.headers {
		fmt.Println("uh.headers", hk, hv)
		req.Header.Set(hk, hv)
	}

	for hk, hv := range additionalHeaders {
		req.Header.Set(hk, hv)
	}
	//godump.Dump(req)
	uh.lastRequestParams = params
	return uh.client.Do(req)
}

func (uh *HttpClient) RequestPlain(method string, path string, additionalHeaders map[string]string) (body []byte, err error) {
	resp, err := uh.Request(method, path, additionalHeaders)
	if err != nil {
		return
	}
	uh.lastRespStatusCode = resp.StatusCode

	//godump.Dump(resp)
	defer resp.Body.Close()

	body, err = io.ReadAll(resp.Body)
	return
}

func (uh *HttpClient) RequestJson(v interface{}, method string, path string, additionalHeaders map[string]string) (err error) {

	body, err := uh.RequestPlain(method, path, additionalHeaders)
	if err != nil {
		return
	}
	//fmt.Println(string(body))
	err = json.Unmarshal(body, &v)
	if err != nil {
		return
	}
	return
}

func (uh *HttpClient) RequestJsonApiAndDecode(v interface{}, method string, path string, additionalHeaders map[string]string) (err error) {

	apiReturn := &ApiReturnJson{}
	err = uh.RequestJson(apiReturn, method, path, additionalHeaders)
	if nil != err {
		return
	}
	if !apiReturn.Status {
		err = fmt.Errorf("code: %d ,message: %s ,errors: %+v ", apiReturn.Code, apiReturn.Message, apiReturn.Errors)
		return
	}
	if enStr, ok := apiReturn.Data.(string); ok {
		err = uh.aesEncryptor.Decrypt(enStr, &v)
		if nil != err {
			return
		}
	} else {
		err = fmt.Errorf("返回数据的data字段不是加密字符串")
	}

	return
}

func (uh *HttpClient) PostFile(path string, filedName string, file string, headers map[string]string) (resp *http.Response, err error) {
	remoteUrl, params, err := uh.BuildRemoteUrlAndParams("post", path)

	if nil != err {
		return
	}
	pipeReader, pipeWriter := io.Pipe()
	multipartWriter := multipart.NewWriter(pipeWriter)

	go func() {
		defer pipeWriter.Close()
		defer multipartWriter.Close()
		for pk, _ := range params {
			multipartWriter.WriteField(pk, params.Get(pk))
		}

		part, err1 := multipartWriter.CreateFormFile(filedName, filepath.Base(file))
		if err1 != nil {
			err = err1
			return
		}
		fileHandle, err1 := os.Open(file)
		if err1 != nil {
			err = err1
			return
		}
		defer fileHandle.Close()
		if _, err1 = io.Copy(part, fileHandle); nil != err1 {
			err = err1
			return
		}
	}()

	req, err := http.NewRequest("POST", remoteUrl, pipeReader)
	if err != nil {
		return
	}
	for hk, hv := range uh.headers {
		fmt.Println("uh.headers", hk, hv)
		req.Header.Set(hk, hv)
	}

	for hk, hv := range headers {
		req.Header.Set(hk, hv)
	}

	req.Header.Set("Content-Type", multipartWriter.FormDataContentType())
	resp, err = uh.client.Do(req)

	return
}
func (uh *HttpClient) PostFilePlain(path string, filedName string, file string, headers map[string]string) (body []byte, err error) {

	resp, err := uh.PostFile(path, filedName, file, headers)
	uh.lastRespStatusCode = resp.StatusCode
	if err != nil {
		return
	}
	//godump.Dump(resp)

	defer resp.Body.Close()

	body, err = io.ReadAll(resp.Body)
	return
}

func (uh *HttpClient) PostFileJson(v interface{}, path string, filedName string, file string, headers map[string]string) error {

	body, err := uh.PostFilePlain(path, filedName, file, headers)
	if err != nil {
		return err
	}
	//fmt.Println(string(body))
	err = json.Unmarshal(body, &v)
	if err != nil {
		return err
	}
	return nil
}

func (uh *HttpClient) DownloadFile(path string, headers map[string]string) (body []byte, contentType string, err error) {
	resp, err := uh.Request("Get", path, headers)
	if err != nil {
		return
	}
	uh.lastRespStatusCode = resp.StatusCode

	defer resp.Body.Close()
	body, err = io.ReadAll(resp.Body)
	if nil != err {
		return
	}
	contentType = resp.Header.Get("Content-Type")

	if resp.StatusCode != 200 {
		body = []byte("")
	}

	return
}

func Request(method string, remoteUrl string, params url.Values, headers map[string]string, timeout ...time.Duration) (resp *http.Response, err error) {

	uh := NewHttpClient("", timeout...)
	uh.WithParams(params)
	return uh.Request(method, remoteUrl, headers)
}

func RequestJson(v interface{}, method string, remoteUrl string, params url.Values, headers map[string]string, timeout ...time.Duration) error {
	uh := NewHttpClient("", timeout...)
	uh.WithParams(params)
	return uh.RequestJson(v, method, remoteUrl, headers)
}

func DownloadFile(remoteUrl string, params url.Values, headers map[string]string, timeout ...time.Duration) (body []byte, contentType string, err error) {
	uh := NewHttpClient("", timeout...)
	uh.WithParams(params)
	return uh.DownloadFile(remoteUrl, headers)
}
