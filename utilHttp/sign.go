package utilHttp

import (
	"fmt"
	"github.com/hilaoyu/go-utils/utilEnc"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

func SignCheckRequest(secret string, r *http.Request, timeValid time.Duration) (err error) {
	timestampStr := r.Form.Get("_timestamp")
	timestamp, _ := strconv.ParseInt(timestampStr, 10, 64)

	//fmt.Println("timestampStr:", timestampStr, ";timestamp:", timestamp, ";min:", time.Now().Add(-1*time.Duration(3)*time.Minute).Unix(), ";max:", time.Now().Add(time.Duration(3)*time.Minute).Unix())

	if timestamp < time.Now().Add(-1*timeValid).Unix() || timestamp > time.Now().Add(timeValid).Unix() {
		err = fmt.Errorf("sign time error")
		return
	}

	sign := r.Form.Get("sign")
	r.Form.Del("sign")
	signStr := r.Form.Encode()
	signStr += secret

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
