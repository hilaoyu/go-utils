package utilXcfsClient

import (
	"fmt"
	"github.com/hilaoyu/go-utils/utilConvert"
	"github.com/hilaoyu/go-utils/utilEnc"
	"github.com/hilaoyu/go-utils/utilHttp"
	"github.com/hilaoyu/go-utils/utilStr"
	"github.com/hilaoyu/go-utils/utilUuid"
	"github.com/hilaoyu/go-utils/utils"
	"net/url"
	"strings"
	"time"
)

const XcfsSchema = "xcfs"
const XcfsUriPrefix = "file/response/"
const XcfsApiUploadTokenGenerate = "/service-api/upload-token/generate"
const XcfsApiUploadTokenVerifyAndUse = "/service-api/upload-tack/verify-and-use"

type XcfsClientConfigApp struct {
	AppId  string `json:"app_id,omitempty"`
	Secret string `json:"secret,omitempty"`
}
type XcfsClientConfigDirectory struct {
	App        string   `json:"app,omitempty"`
	Directory  string   `json:"directory,omitempty"`
	IsPublic   bool     `json:"is_public,omitempty"`
	NamePrefix string   `json:"name_prefix,omitempty"`
	AllowExt   []string `json:"allow_ext,omitempty"`
	MaxSize    int64    `json:"max_size,omitempty"`
}
type XcfsClientConfig struct {
	ServiceUrl  string                                `json:"service_url,omitempty"`
	ApiUrl      string                                `json:"api_url,omitempty"`
	Apps        map[string]*XcfsClientConfigApp       `json:"apps,omitempty"`
	Directories map[string]*XcfsClientConfigDirectory `json:"directories,omitempty"`
}

type XcfsClient struct {
	config *XcfsClientConfig
}
type XcfsUriData struct {
	Path     string
	AppId    string
	FileId   string
	Size     int64
	HashMd5  string
	IsPublic bool
}

type UploadTokenResult struct {
	TokenId      string   `json:"token_id,omitempty"`
	MaxSize      int64    `json:"max_size,omitempty"`
	MaxChunkSize int      `json:"max_chunk_size,omitempty"`
	AllowExt     []string `json:"allow_ext,omitempty"`
	ServiceUrl   string   `json:"service_url,omitempty"`
}

type UriVerifyResult struct {
	FileUri         string `json:"file_uri,omitempty"`
	UploadCompleted bool   `json:"upload_completed,omitempty"`
	Error           string `json:"error,omitempty"`

	Name     string `json:"name,omitempty"`
	IsPublic bool   `json:"is_public,omitempty"`
	FileType string `json:"file_type,omitempty"`
	FileExt  string `json:"file_ext,omitempty"`
	Mimetype string `json:"mimetype,omitempty"`
	Size     int64  `json:"size,omitempty"`
	HashMd5  string `json:"hash_md5,omitempty"`
}

func NewXcfsClient(config *XcfsClientConfig) (xc *XcfsClient) {
	xc = &XcfsClient{config: config}
	return
}

func (xc *XcfsClient) ServiceUrl() string {
	return xc.config.ServiceUrl
}
func (xc *XcfsClient) ApiUrl() string {
	if "" != xc.config.ApiUrl {
		return xc.config.ApiUrl
	}
	return xc.ServiceUrl()
}

func (xc *XcfsClient) IsXcfsUri(uri string) bool {
	u, err := url.Parse(uri)
	if nil != err {
		return false
	}
	s, err := url.Parse(xc.ServiceUrl())
	if nil != err {
		return false
	}

	if XcfsSchema == strings.ToLower(u.Scheme) || strings.ToLower(u.Host) == strings.ToLower(s.Host) {
		return true
	}

	return false
}

func (xc *XcfsClient) GetAppById(appId string) (app *XcfsClientConfigApp) {
	appId = strings.TrimSpace(appId)
	if "" == appId {
		return
	}
	app, _ = utils.MapFind(xc.config.Apps, func(app *XcfsClientConfigApp, s string) bool {
		return app.AppId == appId
	})

	return
}
func (xc *XcfsClient) GetAppByKey(key string) (app *XcfsClientConfigApp) {
	key = strings.TrimSpace(key)
	if "" == key {
		return
	}
	app, _ = xc.config.Apps[key]

	return
}
func (xc *XcfsClient) GetDirectoryByKey(key string) (directory *XcfsClientConfigDirectory) {
	key = strings.TrimSpace(key)
	if "" == key {
		return
	}
	directory, _ = xc.config.Directories[key]

	return
}

func (xc *XcfsClient) ParseUri(uri string) (data *XcfsUriData) {
	uri = strings.TrimSpace(uri)
	if "" == uri || !xc.IsXcfsUri(uri) {
		return
	}

	u, err := url.Parse(uri)
	if nil != err {
		return
	}
	path := strings.TrimSpace(u.Path)
	if "" == path {
		return
	}

	appId := utilStr.Before(utilStr.After(path, XcfsUriPrefix), "/")
	if "" == appId || !utilUuid.IsUuid(appId) {
		return
	}

	fileId := utilStr.Before(utilStr.After(path, appId+"/"), "/")
	fileId = utilStr.Before(fileId, ".")
	if "" == fileId || !utilUuid.IsUuid(fileId) {
		return
	}

	data = &XcfsUriData{
		Path:     path,
		AppId:    appId,
		FileId:   fileId,
		Size:     0,
		HashMd5:  "",
		IsPublic: false,
	}

	if strings.ToLower(u.Scheme) == XcfsSchema && "" != u.Host {
		arr := strings.Split(u.Host, ".")
		if len(arr) >= 1 {
			size, _ := utilConvert.ToInt64(arr[0])
			data.Size = size
		}
		if len(arr) >= 2 {
			data.HashMd5 = arr[1]
		}
		if len(arr) >= 3 && "public" == strings.ToLower(arr[2]) {
			data.IsPublic = true
		}
	}

	return
}

func (xc *XcfsClient) SignResponseUrl(uri string, expiryTimestamp ...int64) (signUrl string) {
	signUrl = uri
	data := xc.ParseUri(uri)
	if nil == data {
		return
	}

	u, err := url.Parse(xc.ServiceUrl())

	if nil != err {
		return
	}

	u.Path = data.Path

	switch data.IsPublic {
	case false:
		if "" == data.AppId {
			break
		}
		app := xc.GetAppById(data.AppId)
		if nil == app {
			break
		}

		signExpiry := int64(0)
		if len(expiryTimestamp) > 0 && expiryTimestamp[0] > 0 {
			signExpiry = expiryTimestamp[0]
		} else {
			signExpiry = time.Now().Add(time.Duration(5) * time.Minute).Unix()
		}

		expiryStr := utilConvert.ToStr(signExpiry)

		sign := utilEnc.Md5(strings.Trim(u.Path, "/") + expiryStr + app.Secret)

		fileUrlQuery := u.Query()
		fileUrlQuery.Set("e", expiryStr)
		fileUrlQuery.Set("s", sign)

		u.RawQuery = fileUrlQuery.Encode()
		break

	default:
		break
	}
	signUrl = u.String()
	return
}

func (xc *XcfsClient) GenerateToken(directory string) (token *UploadTokenResult, err error) {
	directory = strings.TrimSpace(directory)
	if "" == directory {
		err = fmt.Errorf("参数目录为空")
		return
	}

	directoryConf := xc.GetDirectoryByKey(directory)
	if nil == directoryConf {
		err = fmt.Errorf("目录配置不存")
		return
	}

	appConf := xc.GetAppByKey(directoryConf.App)
	if nil == appConf {
		err = fmt.Errorf("应用配置不存在")
		return
	}
	apiUrl := xc.ApiUrl()
	if "" == apiUrl {
		err = fmt.Errorf("服务地址错误")
		return
	}

	apiClient := utilHttp.NewHttpClient(apiUrl)
	apiClient.WithAesEncryptor(appConf.Secret, appConf.AppId)
	apiClient.WithEncryptData(map[string]interface{}{
		"directory":   directoryConf.Directory,
		"is_public":   directoryConf.IsPublic,
		"name_prefix": directoryConf.NamePrefix,
		"allow_ext":   directoryConf.AllowExt,
		"max_size":    directoryConf.MaxSize,
	})

	token = &UploadTokenResult{}
	err = apiClient.RequestJsonApiAndDecrypt(&token, "POST", XcfsApiUploadTokenGenerate, map[string]string{})

	return
}

func (xc *XcfsClient) VerifyAndUse(uris []string) (xcfsUris map[string]string, err error) {
	if len(uris) <= 0 {
		return
	}

	eachAppUris := map[string][]string{}
	xcfsUris = map[string]string{}
	for _, uri := range uris {
		xcfsUris[uri] = uri
		uriData := xc.ParseUri(uri)
		if nil == uriData {
			continue
		}
		_, ok := eachAppUris[uriData.AppId]
		if !ok {
			eachAppUris[uriData.AppId] = []string{uri}
		} else {
			eachAppUris[uriData.AppId] = append(eachAppUris[uriData.AppId], uri)
		}
	}

	if len(eachAppUris) <= 0 {
		return
	}

	apiUrl := xc.ApiUrl()
	if "" == apiUrl {
		err = fmt.Errorf("服务地址错误")
		return
	}
	apiClient := utilHttp.NewHttpClient(apiUrl)

	for appId, appUris := range eachAppUris {
		appConf := xc.GetAppById(appId)
		if nil == appConf {
			err = fmt.Errorf("%v; %s 应用配置不存在", err, appId)
			continue
		}
		apiClient.WithAesEncryptor(appConf.Secret, appConf.AppId)
		apiClient.WithEncryptData(map[string]interface{}{
			"file_uris": appUris,
		})
		uriResult := map[string]*UriVerifyResult{}
		err1 := apiClient.RequestJsonApiAndDecrypt(&uriResult, "POST", XcfsApiUploadTokenVerifyAndUse, map[string]string{})
		if nil != err1 {
			err = fmt.Errorf("%v; %s 接口调用出错: %v", err, appId, err1)
			continue
		}

		for urk, urv := range uriResult {
			if "" != urv.Error {
				err = fmt.Errorf("%v; %s 接口调用出错: %v", err, urk, urv.Error)
				continue
			}
			if "" != urv.FileUri {
				xcfsUris[urk] = urv.FileUri
			}

		}

	}

	return
}
