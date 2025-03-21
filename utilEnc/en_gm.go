package utilEnc

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/hilaoyu/go-utils/utilRandom"
	"github.com/hilaoyu/go-utils/utils"
	"github.com/tjfoc/gmsm/sm2"
	"github.com/tjfoc/gmsm/sm4"
	"github.com/tjfoc/gmsm/x509"
)

type GmEncryptor struct {
	sm2privateKey *sm2.PrivateKey
	sm2publicKey  *sm2.PublicKey

	sm4Key []byte
}

func GmSm2CreateKeys() (privateKey *sm2.PrivateKey, publicKey *sm2.PublicKey, err error) {
	// 生成私钥文件
	privateKey, err = sm2.GenerateKey(utilRandom.RandReader())
	if err != nil {
		return
	}
	// 生成公钥文件
	publicKey = &privateKey.PublicKey
	return
}
func GmSm2CreateKeysPem() (privateKeyPem []byte, publicKeyPem []byte, err error) {

	// 生成私钥文件
	privateKey, publicKey, err := GmSm2CreateKeys()
	if err != nil {
		return
	}
	privateKeyPem, err = x509.WritePrivateKeyToPem(privateKey, nil)
	if err != nil {
		return
	}

	publicKeyPem, err = x509.WritePublicKeyToPem(publicKey)
	if err != nil {
		return
	}

	return
}

func GmSm4CreateKey() []byte {
	return []byte(utilRandom.RandString(sm4.BlockSize))
}

func NewGmEncryptor() (encryptor *GmEncryptor) {
	encryptor = &GmEncryptor{}
	return
}

func (r *GmEncryptor) SetSm2PrivateKey(privateKey []byte, pwd []byte) (key *sm2.PrivateKey, err error) {
	key, err = x509.ReadPrivateKeyFromPem(privateKey, pwd)
	if nil != err {
		return
	}
	r.sm2privateKey = key
	return
}
func (r *GmEncryptor) SetSm2PublicKey(publicKey []byte) (key *sm2.PublicKey, err error) {

	key, err = x509.ReadPublicKeyFromPem(publicKey)
	if nil != err {
		return
	}
	r.sm2publicKey = key
	return
}

func (r *GmEncryptor) Sm2PrivateKeySign(data []byte) (sign []byte, err error) {

	if nil == r.sm2privateKey {
		err = fmt.Errorf("private key is nil")
		return
	}
	sign, err = r.sm2privateKey.Sign(utilRandom.RandReader(), data, nil)
	return
}
func (r *GmEncryptor) Sm2PrivateKeySignAndBase64(data []byte) (sign string, err error) {

	signByte, err := r.Sm2PrivateKeySign(data)
	if err != nil {
		return
	}
	sign = base64.StdEncoding.EncodeToString(signByte)
	return
}

func (r *GmEncryptor) Sm2PublicKeyVerifySign(data []byte, sign []byte) (err error) {

	if nil == r.sm2publicKey {
		err = fmt.Errorf("public key is nil")
		return
	}
	if !r.sm2publicKey.Verify(data, sign) {
		err = fmt.Errorf("sign verify faild")
	}

	return
}
func (r *GmEncryptor) Sm2Base64DecodeAndPublicKeyVerifySign(data []byte, sign string) (err error) {
	decodeSign, err := base64.StdEncoding.DecodeString(sign)
	if err != nil {
		return
	}
	return r.Sm2PublicKeyVerifySign(data, decodeSign)
}

func (r *GmEncryptor) Sm2PublicKeyEncrypt(data []byte) (enData []byte, err error) {
	if nil == r.sm2publicKey {
		err = fmt.Errorf("public key is nil")
		return
	}
	enData, err = sm2.Encrypt(r.sm2publicKey, data, utilRandom.RandReader(), sm2.C1C3C2)
	return
}
func (r *GmEncryptor) Sm2MarshalAndPublicKeyEncrypt(data interface{}) (enData []byte, err error) {
	jsonByte, err := json.Marshal(data)
	if nil != err {
		err = fmt.Errorf("data to json  error: %+v", err)
		return
	}
	enData, err = r.Sm2PublicKeyEncrypt(jsonByte)
	return
}
func (r *GmEncryptor) Sm2PublicKeyEncryptAndBase64(data []byte) (enData string, err error) {
	dataByte, err := r.Sm2PublicKeyEncrypt(data)
	if nil != err {
		return
	}
	enData = base64.StdEncoding.EncodeToString(dataByte)
	return
}
func (r *GmEncryptor) Sm2MarshalAndPublicKeyEncryptAndBase64(data interface{}) (enData string, err error) {
	jsonByte, err := json.Marshal(data)
	if nil != err {
		err = fmt.Errorf("data to json  error: %+v", err)
		return
	}
	enData, err = r.Sm2PublicKeyEncryptAndBase64(jsonByte)
	return
}

func (r *GmEncryptor) Sm2PrivateKeyDecrypt(cipher []byte) (data []byte, err error) {
	if nil == r.sm2privateKey {
		err = fmt.Errorf("private key is nil")
		return
	}
	decrypt, err := sm2.Decrypt(r.sm2privateKey, cipher, sm2.C1C3C2)
	if err != nil {
		return nil, err
	}
	return decrypt, nil
}
func (r *GmEncryptor) Sm2PrivateKeyDecryptAndUnmarshal(cipher []byte, v interface{}) (err error) {
	data, err := r.Sm2PrivateKeyDecrypt(cipher)
	if nil != err {
		return
	}
	err = json.Unmarshal(data, &v)
	return
}
func (r *GmEncryptor) Sm2Base64DecodeAndPrivateKeyDecrypt(cipherText string) (data []byte, err error) {
	decodeCipher, err := base64.StdEncoding.DecodeString(cipherText)
	if nil != err {
		return
	}

	data, err = r.Sm2PrivateKeyDecrypt(decodeCipher)

	return
}
func (r *GmEncryptor) Sm2Base64DecodeAndPrivateKeyDecryptAndUnmarshal(cipherText string, v interface{}) (err error) {
	data, err := r.Sm2Base64DecodeAndPrivateKeyDecrypt(cipherText)
	if nil != err {
		return
	}
	err = json.Unmarshal(data, &v)
	return
}

func (r *GmEncryptor) SetSm4Key(key []byte) (err error) {

	r.sm4Key = key
	return
}

func (r *GmEncryptor) Sm4Encrypt(data []byte) (enData []byte, err error) {
	enData, err = sm4.Sm4Ecb(r.sm4Key, data, true)
	return
}
func (r *GmEncryptor) Sm4MarshalAndEncrypt(data interface{}) (enData []byte, err error) {
	jsonByte, err := json.Marshal(data)
	if nil != err {
		err = fmt.Errorf("data to json  error: %+v", err)
		return
	}
	enData, err = r.Sm4Encrypt(jsonByte)

	return
}
func (r *GmEncryptor) Sm4EncryptAndBase64(data []byte) (enData string, err error) {
	enByte, err := r.Sm4Encrypt(data)
	if nil != err {
		return
	}
	enData = base64.StdEncoding.EncodeToString(enByte)
	return
}
func (r *GmEncryptor) Sm4MarshalAndEncryptAndBase64(data interface{}) (enData string, err error) {
	enByte, err := r.Sm4MarshalAndEncrypt(data)
	if nil != err {
		return
	}
	enData = base64.StdEncoding.EncodeToString(enByte)
	return
}
func (r *GmEncryptor) Sm4Decrypt(enData []byte) (data []byte, err error) {
	data, err = sm4.Sm4Ecb(r.sm4Key, enData, false)
	return
}
func (r *GmEncryptor) Sm4DecryptAndUnmarshal(enData []byte, v interface{}) (err error) {
	data, err := r.Sm4Decrypt(enData)
	if nil != err {
		return
	}
	err = json.Unmarshal(data, &v)
	return
}
func (r *GmEncryptor) Sm4Base64DecodeAndDecrypt(cipherText string) (data []byte, err error) {
	decodeCipher, err := base64.StdEncoding.DecodeString(cipherText)
	if nil != err {
		return
	}

	data, err = r.Sm4Decrypt(decodeCipher)

	return
}
func (r *GmEncryptor) Sm4Base64DecodeAndDecryptAndUnmarshal(cipherText string, v interface{}) (err error) {
	data, err := r.Sm4Base64DecodeAndDecrypt(cipherText)
	if nil != err {
		return
	}
	err = json.Unmarshal(data, &v)
	return
}

func (r *GmEncryptor) ApiDataEncrypt(data interface{}) (enStr string, err error) {
	sm4Key := GmSm4CreateKey()

	sm4r := NewGmEncryptor()
	_ = sm4r.SetSm4Key(sm4Key)
	enData, err := sm4r.Sm4MarshalAndEncrypt(data)
	if nil != err {
		return
	}
	enKey, err := r.Sm2PublicKeyEncrypt(sm4Key)
	if nil != err {
		return
	}

	enStr = string(utils.Base64EncodeFormByte(append(enKey, enData...)))
	return
}

func (r *GmEncryptor) ApiDataDecrypt(enStr string, v interface{}) (err error) {
	enByte, err := base64.StdEncoding.DecodeString(enStr)
	if nil != err {
		err = fmt.Errorf("base64解码错误: %v", err)
		return
	}

	sm2EnDataLen := 97 + sm4.BlockSize

	sm2EnData := enByte[:sm2EnDataLen]
	sm4EnData := enByte[sm2EnDataLen:]

	sm2DeData, err := r.Sm2PrivateKeyDecrypt(sm2EnData)
	if nil != err {
		err = fmt.Errorf("sm2解密失败: %v", err)
		return
	}

	sm4Key := sm2DeData[:sm4.BlockSize]

	sm4r := NewGmEncryptor()
	_ = sm4r.SetSm4Key(sm4Key)

	err = sm4r.Sm4DecryptAndUnmarshal(sm4EnData, v)

	return
}
