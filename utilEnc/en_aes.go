package utilEnc

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
)

type enDataJson struct {
	Iv    string `json:"iv"`
	Value string `json:"value"`
}
type AesEncryptor struct {
	secret string
}

func NewAesEncryptor(secret string) (aesEncryptor *AesEncryptor) {
	return &AesEncryptor{secret: secret}
}
func (ae *AesEncryptor) SetSecret(secret string) *AesEncryptor {
	ae.secret = secret
	return ae
}
func (ae *AesEncryptor) GetSecret() string {
	return ae.secret
}

func (ae *AesEncryptor) Encrypt(data interface{}) (string, error) {
	jsonStr, err := json.Marshal(data)
	if nil != err {
		return "", fmt.Errorf("Aes data to json  error: %+v", err)
	}
	return ae.EncryptString(string(jsonStr))
}
func (ae *AesEncryptor) EncryptString(data string) (string, error) {

	enData, iv, err := AesCBCEncrypt([]byte(data), []byte(ae.secret))
	if nil != err {
		return "", fmt.Errorf("aes AesCBCEncrypt error: %+v", err)
	}

	enJson := enDataJson{
		Iv:    base64.StdEncoding.EncodeToString(iv),
		Value: base64.StdEncoding.EncodeToString(enData),
	}

	enStr, err := json.Marshal(enJson)
	if nil != err {
		return "", fmt.Errorf("aes to return json error: %+v", err)
	}
	return base64.StdEncoding.EncodeToString(enStr), err
}

func (ae *AesEncryptor) Decrypt(data string, v interface{}) (err error) {
	jsonStr, err := ae.DecryptString(data)
	if nil != err {
		return fmt.Errorf("decrypt to string  error: %+v", err)
	}

	return json.Unmarshal([]byte(jsonStr), &v)
}

func (ae *AesEncryptor) DecryptString(data string) (string, error) {

	inJson, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		return "", fmt.Errorf("in data error: %+v", err)
	}

	var inJsonData enDataJson
	err = json.Unmarshal(inJson, &inJsonData)
	if err != nil {
		return "", fmt.Errorf("in data json error: %+v", err)
	}
	//godump.Dump(inJsonData);

	dataByte, err := base64.StdEncoding.DecodeString(inJsonData.Value)
	if err != nil {
		return "", fmt.Errorf("in data value error: %+v", err)
	}
	iv, err := base64.StdEncoding.DecodeString(inJsonData.Iv)
	if err != nil {
		return "", fmt.Errorf("in data iv error: %+v", err)
	}
	out, err := AesCbcDecrypt(dataByte, []byte(ae.secret), iv)
	//fmt.Println(out)
	return string(out), err
}

// aes加密，填充秘钥key的16位，24,32分别对应AES-128, AES-192, or AES-256.
func AesCBCEncrypt(rawData []byte, key []byte) ([]byte, []byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, nil, err
	}

	//填充原文
	blockSize := block.BlockSize()
	rawData = PKCS7Padding(rawData, blockSize)
	//初始向量IV必须是唯一，但不需要保密

	//block大小 16
	iv := make([]byte, blockSize)

	if _, err := rand.Reader.Read(iv); err != nil {
		return nil, nil, err
	}
	//block大小和初始向量大小一定要一致
	encrypted := make([]byte, len(rawData))
	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(encrypted, rawData)

	return encrypted, iv, nil
}

func AesCbcDecrypt(enData []byte, key []byte, iv []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return make([]byte, 0), fmt.Errorf("key length error: %+v", err)
	}
	blockSize := block.BlockSize()

	// CBC mode always works in whole blocks.
	if len(enData) < aes.BlockSize {
		panic("cipher text must be longer than block size")
	} else if len(enData)%blockSize != 0 {
		panic("ciphertext is not a multiple of the block size")
	}
	decrypter := cipher.NewCBCDecrypter(block, iv)
	deData := make([]byte, len(enData))
	decrypter.CryptBlocks(deData, enData)
	deData = PKCS7UnPadding(deData)
	return deData, nil
}
func NullUnPadding(in []byte) []byte {
	return bytes.TrimRight(in, string([]byte{0}))
}

func PKCS7Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

func PKCS7UnPadding(data []byte) []byte {
	length := len(data)
	unpadding := int(data[length-1])
	return data[:(length - unpadding)]
}
