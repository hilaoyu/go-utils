package utilImg

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"github.com/hilaoyu/go-utils/utilFile"
	"github.com/hilaoyu/go-utils/utilHttp"
	"github.com/hilaoyu/go-utils/utils"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"time"
)

func SaveBase64ToFile(data string, dir string, dst string) (string, error) {
	if "" == dir {
		return "", fmt.Errorf("图片保存位置错误")
	}
	if !utilFile.CheckDir(dir) {
		return "", fmt.Errorf("图片保存目录错误")
	}

	b, _ := regexp.MatchString(`^data:\s*image\/(\w+);base64,`, data)
	if !b {
		return "", fmt.Errorf("图片数据错误")
	}
	re, _ := regexp.Compile(`^data:\s*image\/(\w+);base64,`)
	allData := re.FindAllSubmatch([]byte(data), 2)
	fileType := string(allData[0][1]) //png ，jpeg 后缀获取

	base64Str := re.ReplaceAllString(data, "")

	dataByte, _ := base64.StdEncoding.DecodeString(base64Str)

	if "" == dst {
		dst = strconv.FormatInt(time.Now().UnixMicro(), 10) + "." + fileType
	}

	path := filepath.Join(dir, dst)

	err := ioutil.WriteFile(path, dataByte, 0666)
	if nil != err {
		return "", fmt.Errorf("写入文件错误")

	}

	return path, err
}

// Base64FromLocal reads a local file and returns
// the base64 encoded version.
func Base64FromLocal(src string) (string, error) {
	var b bytes.Buffer

	fileExists := utilFile.Exists(src)
	if !fileExists {
		return "", fmt.Errorf("File does not exist\n")
	}

	file, err := os.Open(src)
	if err != nil {
		return "", fmt.Errorf("Error opening file\n")
	}

	_, err = b.ReadFrom(file)
	if err != nil {
		return "", fmt.Errorf("Error reading file to buffer\n")
	}

	return Base64FromBuffer(b), nil
}

// Base64FromRemote is a better named function that
// presently calls NewImage which will be deprecated.
// Function accepts an RFC compliant URL and returns
// a base64 encoded result.
func Base64FromRemote(remoteUrl string) (out string, err error) {

	image, mime, err := utilHttp.DownloadFile(remoteUrl, url.Values{}, map[string]string{})
	if nil != err {
		return
	}
	enc := utils.Base64EncodeFormByte(image)

	out = format(enc, mime)
	return
}

// Base64FromBuffer accepts a buffer and returns a
// base64 encoded string.
func Base64FromBuffer(buf bytes.Buffer) string {
	enc := utils.Base64EncodeFormByte(buf.Bytes())
	mime := http.DetectContentType(buf.Bytes())

	return format(enc, mime)
}

// format is an abstraction of the mime switch to create the
// acceptable base64 string needed for browsers.
func format(enc []byte, mime string) string {
	switch mime {
	case "image/gif", "image/jpeg", "image/pjpeg", "image/png", "image/tiff":
		return fmt.Sprintf("data:%s;base64,%s", mime, enc)
	default:
	}

	return fmt.Sprintf("data:image/png;base64,%s", enc)
}
