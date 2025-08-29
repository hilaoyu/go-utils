package utilEnc

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/hilaoyu/go-utils/utilRandom"
	"github.com/tjfoc/gmsm/sm4"
)

type GmSm4Encryptor struct {
	sm4Key []byte
}

func GmSm4CreateKey() []byte {
	return []byte(utilRandom.RandString(sm4.BlockSize))
}

func NewGmSm4Encryptor(key []byte) (encryptor *GmSm4Encryptor) {
	encryptor = &GmSm4Encryptor{sm4Key: key}
	return
}

func (r *GmSm4Encryptor) SetSm4Key(key []byte) (err error) {

	r.sm4Key = key
	return
}

func (r *GmSm4Encryptor) Sm4Encrypt(data []byte) (enData []byte, err error) {
	enData, err = sm4.Sm4Ecb(r.sm4Key, data, true)
	return
}
func (r *GmSm4Encryptor) Sm4MarshalAndEncrypt(data interface{}) (enData []byte, err error) {
	jsonByte, err := json.Marshal(data)
	if nil != err {
		err = fmt.Errorf("data to json  error: %+v", err)
		return
	}
	enData, err = r.Sm4Encrypt(jsonByte)

	return
}
func (r *GmSm4Encryptor) Sm4EncryptAndBase64(data []byte) (enData string, err error) {
	enByte, err := r.Sm4Encrypt(data)
	if nil != err {
		return
	}
	enData = base64.StdEncoding.EncodeToString(enByte)
	return
}
func (r *GmSm4Encryptor) Sm4MarshalAndEncryptAndBase64(data interface{}) (enData string, err error) {
	enByte, err := r.Sm4MarshalAndEncrypt(data)
	if nil != err {
		return
	}
	enData = base64.StdEncoding.EncodeToString(enByte)
	return
}
func (r *GmSm4Encryptor) Sm4Decrypt(enData []byte) (data []byte, err error) {
	data, err = sm4.Sm4Ecb(r.sm4Key, enData, false)
	return
}
func (r *GmSm4Encryptor) Sm4DecryptAndUnmarshal(enData []byte, v interface{}) (err error) {
	data, err := r.Sm4Decrypt(enData)
	if nil != err {
		return
	}
	err = json.Unmarshal(data, &v)
	return
}
func (r *GmSm4Encryptor) Sm4Base64DecodeAndDecrypt(cipherText string) (data []byte, err error) {
	decodeCipher, err := base64.StdEncoding.DecodeString(cipherText)
	if nil != err {
		return
	}

	data, err = r.Sm4Decrypt(decodeCipher)

	return
}
func (r *GmSm4Encryptor) Sm4Base64DecodeAndDecryptAndUnmarshal(cipherText string, v interface{}) (err error) {
	data, err := r.Sm4Base64DecodeAndDecrypt(cipherText)
	if nil != err {
		return
	}
	err = json.Unmarshal(data, &v)
	return
}

func (r *GmSm4Encryptor) EncryptorType() string {
	return ApiDataEncryptorTypeGmSm4
}

func (r *GmSm4Encryptor) ApiDataEncrypt(data interface{}) (enStr string, err error) {

	enStr, err = r.Sm4MarshalAndEncryptAndBase64(data)
	return
}

func (r *GmSm4Encryptor) ApiDataDecrypt(enStr string, v interface{}) (err error) {

	err = r.Sm4Base64DecodeAndDecryptAndUnmarshal(enStr, v)

	return
}
