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
	secret      string
	iv          []byte
	cipherBlock cipher.Block
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
func (ae *AesEncryptor) GetCipherBlock() (block cipher.Block, err error) {
	if nil == ae.cipherBlock {
		ae.cipherBlock, err = aes.NewCipher([]byte(ae.secret))
		if nil != err {
			return
		}
	}
	block = ae.cipherBlock
	return
}
func (ae *AesEncryptor) GetBlockSize() (size int, err error) {
	block, err := ae.GetCipherBlock()
	if nil != err {
		return
	}
	size = block.BlockSize()
	return
}
func (ae *AesEncryptor) RandIv() (iv []byte, err error) {
	length, err := ae.GetBlockSize()
	if nil != err {
		return
	}
	iv = make([]byte, length)
	_, err = rand.Reader.Read(iv)
	return
}

// aes加密，填充秘钥key的16位，24,32分别对应AES-128, AES-192, or AES-256.
func (ae *AesEncryptor) EncryptByte(rawData []byte, iv []byte) (enData []byte, err error) {
	block, err := ae.GetCipherBlock()
	if nil != err {
		return
	}

	/*if nil == iv {
		iv, err = ae.RandIv()
		if nil != err {
			return
		}
	}*/
	rawData = PKCS7Padding(rawData, block.BlockSize())
	enData = make([]byte, len(rawData))
	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(enData, rawData)

	return
}

func (ae *AesEncryptor) Encrypt(data interface{}) (string, error) {
	jsonStr, err := json.Marshal(data)
	if nil != err {
		return "", fmt.Errorf("Aes data to json  error: %+v", err)
	}
	return ae.EncryptString(string(jsonStr))
}
func (ae *AesEncryptor) EncryptString(data string) (enStr string, err error) {

	iv, err := ae.RandIv()
	if nil != err {
		return "", err
	}
	enData, err := ae.EncryptByte([]byte(data), iv)
	if nil != err {
		return "", fmt.Errorf("aes AesCBCEncrypt error: %+v", err)
	}

	enJson := enDataJson{
		Iv:    base64.StdEncoding.EncodeToString(iv),
		Value: base64.StdEncoding.EncodeToString(enData),
	}

	jsonByte, err := json.Marshal(enJson)
	if nil != err {
		err = fmt.Errorf("aes to return json error: %+v", err)
		return
	}
	enStr = base64.StdEncoding.EncodeToString(jsonByte)
	return
}

func (ae *AesEncryptor) DecryptByte(enData []byte, iv []byte) (deData []byte, err error) {
	block, err := ae.GetCipherBlock()
	if nil != err {
		return
	}
	blockSize := block.BlockSize()
	// CBC mode always works in whole blocks.
	if len(enData) < aes.BlockSize {
		err = fmt.Errorf("cipher text must be longer than block size")
		return
	} else if len(enData)%blockSize != 0 {
		err = fmt.Errorf("ciphertext is not a multiple of the block size")
		return
	}
	decryptor := cipher.NewCBCDecrypter(block, iv)
	deData = make([]byte, len(enData))
	decryptor.CryptBlocks(deData, enData)
	deData, err = PKCS7UnPadding(deData)

	return
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

	//fmt.Println(out)
	out, err := ae.DecryptByte(dataByte, iv)
	return string(out), err
}

func PKCS7Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

func PKCS7UnPadding(data []byte) (b []byte, err error) {
	length := len(data)
	unpadding := int(data[length-1])
	end := length - unpadding
	if end < 0 {
		err = fmt.Errorf("unpadding index error")
		return
	}
	b = data[:end]

	return
}
