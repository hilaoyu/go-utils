package utilEnc

const ApiDataEncryptorTypeRsa = " API_DATA_ENCRYPTOR_RSA"
const ApiDataEncryptorTypeAes = " API_DATA_ENCRYPTOR_AES"

type ApiDataEncryptor interface {
	ApiDataEncrypt(data interface{}) (enStr string, err error)
	ApiDataDecrypt(enStr string, v interface{}) (err error)
}
