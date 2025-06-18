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

type GmSm2Encryptor struct {
	sm2privateKey *sm2.PrivateKey
	sm2publicKey  *sm2.PublicKey
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

func NewGmSm2Encryptor() (encryptor *GmSm2Encryptor) {
	encryptor = &GmSm2Encryptor{}
	return
}

func (r *GmSm2Encryptor) SetSm2PrivateKey(privateKey []byte, pwd []byte) (key *sm2.PrivateKey, err error) {
	key, err = x509.ReadPrivateKeyFromPem(privateKey, pwd)
	if nil != err {
		return
	}
	r.sm2privateKey = key
	return
}
func (r *GmSm2Encryptor) SetSm2PublicKey(publicKey []byte) (key *sm2.PublicKey, err error) {

	key, err = x509.ReadPublicKeyFromPem(publicKey)
	if nil != err {
		return
	}
	r.sm2publicKey = key
	return
}

func (r *GmSm2Encryptor) Sm2PrivateKeySign(data []byte) (sign []byte, err error) {

	if nil == r.sm2privateKey {
		err = fmt.Errorf("private key is nil")
		return
	}
	sign, err = r.sm2privateKey.Sign(utilRandom.RandReader(), data, nil)
	return
}
func (r *GmSm2Encryptor) Sm2PrivateKeySignAndBase64(data []byte) (sign string, err error) {

	signByte, err := r.Sm2PrivateKeySign(data)
	if err != nil {
		return
	}
	sign = base64.StdEncoding.EncodeToString(signByte)
	return
}

func (r *GmSm2Encryptor) Sm2PublicKeyVerifySign(data []byte, sign []byte) (err error) {

	if nil == r.sm2publicKey {
		err = fmt.Errorf("public key is nil")
		return
	}
	if !r.sm2publicKey.Verify(data, sign) {
		err = fmt.Errorf("sign verify faild")
	}

	return
}
func (r *GmSm2Encryptor) Sm2Base64DecodeAndPublicKeyVerifySign(data []byte, sign string) (err error) {
	decodeSign, err := base64.StdEncoding.DecodeString(sign)
	if err != nil {
		return
	}
	return r.Sm2PublicKeyVerifySign(data, decodeSign)
}

func (r *GmSm2Encryptor) Sm2PublicKeyEncrypt(data []byte) (enData []byte, err error) {
	if nil == r.sm2publicKey {
		err = fmt.Errorf("public key is nil")
		return
	}
	enData, err = sm2.Encrypt(r.sm2publicKey, data, utilRandom.RandReader(), sm2.C1C3C2)
	return
}
func (r *GmSm2Encryptor) Sm2MarshalAndPublicKeyEncrypt(data interface{}) (enData []byte, err error) {
	jsonByte, err := json.Marshal(data)
	if nil != err {
		err = fmt.Errorf("data to json  error: %+v", err)
		return
	}
	enData, err = r.Sm2PublicKeyEncrypt(jsonByte)
	return
}
func (r *GmSm2Encryptor) Sm2PublicKeyEncryptAndBase64(data []byte) (enData string, err error) {
	dataByte, err := r.Sm2PublicKeyEncrypt(data)
	if nil != err {
		return
	}
	enData = base64.StdEncoding.EncodeToString(dataByte)
	return
}
func (r *GmSm2Encryptor) Sm2MarshalAndPublicKeyEncryptAndBase64(data interface{}) (enData string, err error) {
	jsonByte, err := json.Marshal(data)
	if nil != err {
		err = fmt.Errorf("data to json  error: %+v", err)
		return
	}
	enData, err = r.Sm2PublicKeyEncryptAndBase64(jsonByte)
	return
}

func (r *GmSm2Encryptor) Sm2PrivateKeyDecrypt(cipher []byte) (data []byte, err error) {
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
func (r *GmSm2Encryptor) Sm2PrivateKeyDecryptAndUnmarshal(cipher []byte, v interface{}) (err error) {
	data, err := r.Sm2PrivateKeyDecrypt(cipher)
	if nil != err {
		return
	}
	err = json.Unmarshal(data, &v)
	return
}
func (r *GmSm2Encryptor) Sm2Base64DecodeAndPrivateKeyDecrypt(cipherText string) (data []byte, err error) {
	decodeCipher, err := base64.StdEncoding.DecodeString(cipherText)
	if nil != err {
		return
	}

	data, err = r.Sm2PrivateKeyDecrypt(decodeCipher)

	return
}
func (r *GmSm2Encryptor) Sm2Base64DecodeAndPrivateKeyDecryptAndUnmarshal(cipherText string, v interface{}) (err error) {
	data, err := r.Sm2Base64DecodeAndPrivateKeyDecrypt(cipherText)
	if nil != err {
		return
	}
	err = json.Unmarshal(data, &v)
	return
}

func (r *GmSm2Encryptor) EncryptorType() string {
	return ApiDataEncryptorTypeGmSm2
}
func (r *GmSm2Encryptor) ApiDataEncrypt(data interface{}) (enStr string, err error) {
	sm4Key := GmSm4CreateKey()

	sm4r := NewGmSm4Encryptor(sm4Key)
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

func (r *GmSm2Encryptor) ApiDataDecrypt(enStr string, v interface{}) (err error) {
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

	sm4r := NewGmSm4Encryptor(sm4Key)

	err = sm4r.Sm4DecryptAndUnmarshal(sm4EnData, v)

	return
}
