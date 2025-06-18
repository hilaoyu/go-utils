package utilEnc

const ApiDataEncryptorTypeRsa = " API_DATA_ENCRYPTOR_RSA"
const ApiDataEncryptorTypeAes = " API_DATA_ENCRYPTOR_AES"
const ApiDataEncryptorTypeGmSm2 = " API_DATA_ENCRYPTOR_GM_SM2"
const ApiDataEncryptorTypeGmSm4 = " API_DATA_ENCRYPTOR_GM_SM4"

type ApiDataEncryptor interface {
	EncryptorType() string
	ApiDataEncrypt(data interface{}) (enStr string, err error)
	ApiDataDecrypt(enStr string, v interface{}) (err error)
}
