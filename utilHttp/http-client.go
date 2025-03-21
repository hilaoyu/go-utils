package utilHttp

import (
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/hilaoyu/go-utils/utilEnc"
	"github.com/hilaoyu/go-utils/utilLogger"
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
		params:                  url.Values{},
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
	uh.encryptor = utilEnc.NewAesEncryptor(secret)
	uh.encryptorType = utilEnc.ApiDataEncryptorTypeAes
	uh.aesEncAppId = appId
	return uh.buildClient()
}
func (uh *HttpClient) WithRsaEncryptor(publicKey []byte, privateKey []byte, appId string) *HttpClient {
	encryptor := utilEnc.NewRsaEncryptor()
	var err error
	if len(publicKey) > 0 {
		_, err = encryptor.SetPublicKey(publicKey)

		if nil != err {
			uh.logError(fmt.Sprintf("http client: %v", err))
		}
	}
	if len(privateKey) > 0 {
		_, err = encryptor.SetPrivateKey(privateKey)
		if nil != err {
			uh.logError(fmt.Sprintf("http client: %v", err))
		}
	}
	uh.encryptor = encryptor
	uh.encryptorType = utilEnc.ApiDataEncryptorTypeRsa
	uh.aesEncAppId = appId
	return uh.buildClient()
}
func (uh *HttpClient) WithGmEncryptor(publicKey []byte, privateKey []byte, appId string) *HttpClient {
	encryptor := utilEnc.NewGmEncryptor()
	var err error
	if len(publicKey) > 0 {
		_, err = encryptor.SetSm2PublicKey(publicKey)

		if nil != err {
			uh.logError(fmt.Sprintf("http client: %v", err))
		}
	}
	if len(privateKey) > 0 {
		_, err = encryptor.SetSm2PrivateKey(privateKey, nil)
		if nil != err {
			uh.logError(fmt.Sprintf("http client: %v", err))
		}
	}
	uh.encryptor = encryptor
	uh.encryptorType = utilEnc.ApiDataEncryptorTypeGm
	uh.aesEncAppId = appId
	return uh.buildClient()
}
func (uh *HttpClient) GetEncryptor() utilEnc.ApiDataEncryptor {
	return uh.encryptor
}
func (uh *HttpClient) GetEncryptorType() string {
	return uh.encryptorType
}
func (uh *HttpClient) GetAesEncryptor() (aesEncryptor *utilEnc.AesEncryptor) {
	if nil == uh.encryptor {
		return
	}
	aesEncryptor, ok := uh.encryptor.(*utilEnc.AesEncryptor)
	if !ok {
		aesEncryptor = nil
	}
	return
}
func (uh *HttpClient) GetRsaEncryptor() (aesEncryptor *utilEnc.RsaEncryptor) {
	if nil == uh.encryptor {
		return
	}
	aesEncryptor, ok := uh.encryptor.(*utilEnc.RsaEncryptor)
	if !ok {
		aesEncryptor = nil
	}
	return
}

func (uh *HttpClient) GetGmEncryptor() (aesEncryptor *utilEnc.GmEncryptor) {
	if nil == uh.encryptor {
		return
	}
	aesEncryptor, ok := uh.encryptor.(*utilEnc.GmEncryptor)
	if !ok {
		aesEncryptor = nil
	}
	return
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
			uh.logError(fmt.Sprintf("tls.X509KeyPair err: %v", err))
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
	return uh
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
	return uh
}

func (uh *HttpClient) BuildRemoteUrlAndParams(method string, path string) (remoteUrl string, params url.Values, err error) {
	method = strings.ToUpper(method)
	remoteUrl = path
	if "" != uh.baseUrl {
		remoteUrl = strings.TrimRight(uh.baseUrl, "/") + "/" + strings.TrimLeft(path, "/")
	}

	params = uh.params
	if nil != uh.encryptor {

		if nil != uh.needEncData {
			needEncData := uh.needEncData
			needEncData["_timestamp"] = time.Now().UTC().Unix()
			needEncData["_data_id"] = strconv.FormatInt(time.Now().UTC().UnixNano(), 10)
			enData, err1 := uh.encryptor.ApiDataEncrypt(needEncData)
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
	return uh
}
func (uh *HttpClient) ClearHeader() *HttpClient {
	uh.headers = map[string]string{}
	return uh
}

func (uh *HttpClient) WithParams(params url.Values) *HttpClient {
	uh.params = params
	return uh
}
func (uh *HttpClient) ClearParams() *HttpClient {
	uh.params = url.Values{}
	return uh
}
func (uh *HttpClient) WithJsonData(data interface{}) *HttpClient {
	jsonByte, err := json.Marshal(data)
	if nil == err {
		uh.WithRawBody(string(jsonByte))
	}
	uh.AddHeader("Content-Type", "application/json")
	return uh
}
func (uh *HttpClient) WithRawBody(body string) *HttpClient {
	uh.rawBody = body
	return uh
}
func (uh *HttpClient) ClearRawBody() *HttpClient {
	uh.rawBody = ""
	return uh
}

func (uh *HttpClient) WithEncryptData(data map[string]interface{}) *HttpClient {
	uh.needEncData = data
	return uh
}
func (uh *HttpClient) ClearEncryptData() *HttpClient {
	uh.needEncData = map[string]interface{}{}
	return uh
}

func (uh *HttpClient) GetLastRequestUrl() string {
	return uh.lastRequestUrl
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

	fmt.Println(remoteUrl, params.Encode())

	var body *strings.Reader
	if "" != uh.rawBody {
		body = strings.NewReader(uh.rawBody)
	} else {
		body = strings.NewReader(params.Encode())
	}
	req, err := http.NewRequest(method, remoteUrl, body)

	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	for hk, hv := range uh.headers {
		req.Header.Set(hk, hv)
	}

	for hk, hv := range additionalHeaders {
		req.Header.Set(hk, hv)
	}
	//godump.Dump(req)
	uh.lastRequestUrl = remoteUrl
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

	if nil == additionalHeaders {
		additionalHeaders = map[string]string{}
	}
	additionalHeaders["X-Requested-With"] = "XMLHttpRequest"
	body, err := uh.RequestPlain(method, path, additionalHeaders)
	if err != nil {
		return
	}
	err = json.Unmarshal(body, &v)
	if err != nil {
		return
	}
	return
}

func (uh *HttpClient) RequestJsonApi(v interface{}, method string, path string, additionalHeaders map[string]string) (err error) {

	apiReturn := &ApiDataJson{}
	if nil != v {
		apiReturn.Data = v
	}
	err = uh.RequestJson(apiReturn, method, path, additionalHeaders)
	if nil != err {
		return
	}
	if !apiReturn.Status {
		err = fmt.Errorf("code: %d ,message: %s ,errors: %+v ", apiReturn.Code, apiReturn.Message, apiReturn.Errors)
		return
	}

	return
}

func (uh *HttpClient) RequestJsonApiAndDecrypt(v interface{}, method string, path string, additionalHeaders map[string]string) (err error) {

	apiReturn := &ApiDataJson{}
	err = uh.RequestJson(apiReturn, method, path, additionalHeaders)
	if nil != err {
		return
	}
	if !apiReturn.Status {
		err = fmt.Errorf("code: %d ,message: %s ,errors: %+v ", apiReturn.Code, apiReturn.Message, apiReturn.Errors)
		return
	}
	if nil != v {
		if enStr, ok := apiReturn.Data.(string); ok {
			if nil != uh.encryptor {
				err = uh.encryptor.ApiDataDecrypt(enStr, &v)
			} else {
				err = fmt.Errorf("加密器为空")
			}

			if nil != err {
				return
			}
		} else {
			err = fmt.Errorf("返回数据的data字段不是加密字符串")
		}

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

func (uh *HttpClient) SignRequest(secret string, method string, path string, params url.Values, additionalHeaders map[string]string) (resp *http.Response, err error) {
	params = SignRequestParams(secret, params)
	resp, err = uh.WithParams(params).Request(method, path, additionalHeaders)

	return
}

func (uh *HttpClient) SetLogger(logger *utilLogger.Logger) *HttpClient {
	uh.logger = logger
	return uh
}
func (uh *HttpClient) logInfo(msg interface{}) {
	if nil == uh.logger {
		return
	}

	uh.logger.Info(msg)
	return
}
func (uh *HttpClient) logError(msg interface{}) {
	if nil == uh.logger {
		return
	}

	uh.logger.Error(msg)
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
