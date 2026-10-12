package utilHttp

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/hilaoyu/go-utils/utilCache"
	"github.com/hilaoyu/go-utils/utilEnc"
)

var (
	utilHttpSignDataIdCache = utilCache.NewCache("_util_http_sign_data_id_", time.Duration(5)*time.Minute).RegisterStoreMemory(10000)
)

func SignCheckRequest(secret string, r *http.Request, timeValid time.Duration, debug ...bool) (err error) {
	values, err := GetRequestValues(r)
	if nil != err {
		return
	}

	timestampStr := values.Get("_timestamp")

	timestamp, _ := strconv.ParseInt(timestampStr, 10, 64)

	sign := values.Get("sign")
	dataId := values.Get("_data_id")

	values.Del("sign")
	signStr := values.Encode()
	signStr += secret

	//fmt.Println("timestampStr:", timestampStr, ";timestamp:", timestamp, ";min:", time.Now().Add(-1*time.Duration(3)*time.Minute).Unix(), ";max:", time.Now().Add(time.Duration(3)*time.Minute).Unix())

	if timestamp < time.Now().Add(-1*timeValid).Unix() || timestamp > time.Now().Add(timeValid).Unix() {
		err = fmt.Errorf("sign time error")
		return
	}

	if len(debug) <= 1 || !debug[1] {
		if "" == dataId {
			err = fmt.Errorf("数据Id错误")
			return
		}

		if utilHttpSignDataIdCache.GetBool(dataId) {
			err = fmt.Errorf("数据已处理过了")
			return
		}
		_ = utilHttpSignDataIdCache.Set(dataId, true)
	}

	sign1 := utilEnc.Md5(signStr)

	if sign != sign1 {
		err = fmt.Errorf("sign enc error")
	}
	return
}

func SignRequestParams(secret string, params url.Values) url.Values {
	params.Del("sign")
	params.Set("_timestamp", strconv.FormatInt(time.Now().Unix(), 10))
	params.Set("_data_id", strconv.FormatInt(time.Now().UnixNano(), 10))
	sign := utilEnc.Md5(params.Encode() + secret)
	params.Set("sign", sign)
	return params
}
