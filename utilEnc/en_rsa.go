package utilEnc

import (
	"bytes"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"github.com/hilaoyu/go-utils/utilRandom"
	"github.com/hilaoyu/go-utils/utilSsl"
	"github.com/hilaoyu/go-utils/utils"
)

type RsaEncryptor struct {
	privateKey *rsa.PrivateKey
	publicKey  *rsa.PublicKey
}

func RsaCreateKeys(keyLength int) (privateKey *rsa.PrivateKey, publicKey *rsa.PublicKey, err error) {
	// 生成私钥文件
	privateKey, err = rsa.GenerateKey(rand.Reader, keyLength)
	if err != nil {
		return
	}

	// 生成公钥文件
	publicKey = &privateKey.PublicKey

	return
}
func RsaCreateKeysPem(keyLength int) (privateKeyPem []byte, publicKeyPem []byte, err error) {
	var publicKeyWriter *bytes.Buffer = bytes.NewBufferString("")
	var privateKeyWriter *bytes.Buffer = bytes.NewBufferString("")
	defer func() {
		publicKeyWriter = nil
		privateKeyWriter = nil
	}()
	// 生成私钥文件
	privateKey, publicKey, err := RsaCreateKeys(keyLength)
	if err != nil {
		return
	}
	derStream, err := x509.MarshalPKCS8PrivateKey(privateKey)
	if err != nil {
		return
	}
	block := &pem.Block{
		Type:  "PRIVATE KEY",
		Bytes: derStream,
	}
	err = pem.Encode(privateKeyWriter, block)
	if err != nil {
		return
	}
	privateKeyPem = privateKeyWriter.Bytes()

	derPkix, err := x509.MarshalPKIXPublicKey(publicKey)
	if err != nil {
		return
	}
	block = &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: derPkix,
	}
	err = pem.Encode(publicKeyWriter, block)
	if err != nil {
		return
	}
	publicKeyPem = publicKeyWriter.Bytes()

	return
}

func NewRsaEncryptor() (encryptor *RsaEncryptor) {
	encryptor = &RsaEncryptor{}
	return
}

func (r *RsaEncryptor) SetPrivateKey(privateKey []byte) (key *rsa.PrivateKey, err error) {
	key, err = utilSsl.ParseX509PrivateKeyContent(privateKey)
	if nil != err {
		return
	}
	r.privateKey = key
	return
}

func (r *RsaEncryptor) SetPublicKey(publicKey []byte) (key *rsa.PublicKey, err error) {

	key, err = utilSsl.ParseX509PublicKeyContent(publicKey)
	if nil != err {
		return
	}
	r.publicKey = key
	return
}

func (r *RsaEncryptor) RsaPrivateKeySign(data []byte) (sign []byte, err error) {

	if nil == r.privateKey {
		err = fmt.Errorf("private key is nil")
		return
	}
	h := sha256.New()
	h.Write(data)
	hashed := h.Sum(nil)

	sign, err = rsa.SignPKCS1v15(nil, r.privateKey, crypto.SHA256, hashed)

	return
}
func (r *RsaEncryptor) RsaPrivateKeySignAndBase64(data []byte) (sign string, err error) {

	signByte, err := r.RsaPrivateKeySign(data)
	if err != nil {
		return
	}
	sign = base64.StdEncoding.EncodeToString(signByte)
	return
}

func (r *RsaEncryptor) RsaPublicKeyVerifySign(data []byte, sign []byte) (err error) {

	if nil == r.publicKey {
		err = fmt.Errorf("public key is nil")
		return
	}
	h := sha1.New()
	h.Write(data)
	hashed := h.Sum(nil)
	err = rsa.VerifyPKCS1v15(r.publicKey, crypto.SHA1, hashed, sign)

	return
}
func (r *RsaEncryptor) RsaBase64DecodeAndPublicKeyVerifySign(data []byte, sign string) (err error) {
	decodeSign, err := base64.StdEncoding.DecodeString(sign)
	if err != nil {
		return
	}
	return r.RsaPublicKeyVerifySign(data, decodeSign)
}

func (r *RsaEncryptor) RsaPublicKeyEncrypt(src []byte) (data []byte, err error) {
	if nil == r.publicKey {
		err = fmt.Errorf("public key is nil")
		return
	}
	data, err = rsa.EncryptPKCS1v15(rand.Reader, r.publicKey, src)
	return
}
func (r *RsaEncryptor) RsaMarshalAndPublicKeyEncrypt(data interface{}) (enData []byte, err error) {
	jsonByte, err := json.Marshal(data)
	if nil != err {
		err = fmt.Errorf("rsa data to json  error: %+v", err)
		return
	}
	enData, err = r.RsaPublicKeyEncrypt(jsonByte)
	return
}
func (r *RsaEncryptor) RsaPublicKeyEncryptAndBase64(src []byte) (data string, err error) {
	dataByte, err := r.RsaPublicKeyEncrypt(src)
	if nil != err {
		return
	}
	data = base64.StdEncoding.EncodeToString(dataByte)
	return
}
func (r *RsaEncryptor) RsaMarshalAndPublicKeyEncryptAndBase64(data interface{}) (enData string, err error) {
	jsonByte, err := json.Marshal(data)
	if nil != err {
		err = fmt.Errorf("rsa data to json  error: %+v", err)
		return
	}
	enData, err = r.RsaPublicKeyEncryptAndBase64(jsonByte)
	return
}

func (r *RsaEncryptor) RsaPrivateKeyDecrypt(cipher []byte) (data []byte, err error) {
	if nil == r.privateKey {
		err = fmt.Errorf("private key is nil")
		return
	}
	decrypt, err := rsa.DecryptPKCS1v15(rand.Reader, r.privateKey, cipher)
	if err != nil {
		return []byte{}, err
	}
	return decrypt, nil
}
func (r *RsaEncryptor) RsaPrivateKeyDecryptAndUnmarshal(cipher []byte, v interface{}) (err error) {
	data, err := r.RsaPrivateKeyDecrypt(cipher)
	if nil != err {
		return
	}
	err = json.Unmarshal(data, &v)
	return
}
func (r *RsaEncryptor) RsaBase64DecodeAndPrivateKeyDecrypt(cipherText string) (data []byte, err error) {
	decodeCipher, err := base64.StdEncoding.DecodeString(cipherText)
	if nil != err {
		return
	}

	data, err = r.RsaPrivateKeyDecrypt(decodeCipher)

	return
}
func (r *RsaEncryptor) RsaBase64DecodeAndPrivateKeyDecryptAndUnmarshal(cipherText string, v interface{}) (err error) {
	data, err := r.RsaBase64DecodeAndPrivateKeyDecrypt(cipherText)
	if nil != err {
		return
	}
	err = json.Unmarshal(data, &v)
	return
}

func (r *RsaEncryptor) ApiDataEncryptWithAes(data []byte) (enStr string, err error) {

	aesKey := utilRandom.RandString(16)
	aesEnc := NewAesEncryptor(aesKey)
	iv, err := aesEnc.RandIv()
	if nil != err {
		return
	}
	enData, err := aesEnc.EncryptByte(data, iv)
	if nil != err {
		return
	}
	enKey, err := r.RsaPublicKeyEncrypt(append([]byte(aesKey), iv...))
	if nil != err {
		return
	}
	enStr = string(utils.Base64EncodeFormByte(append(enKey, enData...)))
	return
}
func (r *RsaEncryptor) ApiDataMarshalAndEncryptWithAes(data interface{}) (enStr string, err error) {
	jsonByte, err := json.Marshal(data)
	if nil != err {
		return
	}
	enStr, err = r.ApiDataEncryptWithAes(jsonByte)
	return
}

func (r *RsaEncryptor) ApiDataDecryptWithAes(enStr string) (data []byte, err error) {
	enByte, err := base64.StdEncoding.DecodeString(enStr)
	if nil != err {
		err = fmt.Errorf("base64解码错误: %v", err)
		return
	}

	rsaEnData := enByte[:r.privateKey.Size()]
	aesEnData := enByte[r.privateKey.Size():]

	rsaDeData, err := r.RsaPrivateKeyDecrypt(rsaEnData)
	if nil != err {
		err = fmt.Errorf("rsa解密失败: %v", err)
		return
	}

	aesKey := rsaDeData[:16]
	aesIv := rsaDeData[16:]

	aesEnc := NewAesEncryptor(string(aesKey))
	data, err = aesEnc.DecryptByte(aesEnData, aesIv)
	return
}
func (r *RsaEncryptor) ApiDataDecryptWithAesAndUnmarshal(enStr string, v interface{}) (err error) {
	data, err := r.ApiDataDecryptWithAes(enStr)
	if nil != err {
		return
	}
	err = json.Unmarshal(data, v)
	return
}

func (r *RsaEncryptor) EncryptorType() string {
	return ApiDataEncryptorTypeGmSm4
}
func (r *RsaEncryptor) ApiDataEncrypt(data interface{}) (enStr string, err error) {
	return r.ApiDataMarshalAndEncryptWithAes(data)
}
func (r *RsaEncryptor) ApiDataDecrypt(enStr string, v interface{}) (err error) {
	return r.ApiDataDecryptWithAesAndUnmarshal(enStr, v)
}
