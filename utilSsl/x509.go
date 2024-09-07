package utilSsl

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"
)

func CreateRsaKeys(keyLength int) (privateKey *rsa.PrivateKey, publicKey *rsa.PublicKey, err error) {
	// 生成私钥文件
	privateKey, err = rsa.GenerateKey(rand.Reader, keyLength)
	if err != nil {
		return
	}

	// 生成公钥文件
	publicKey = &privateKey.PublicKey

	return
}
func CreateRsaKeysPem(keyLength int) (privateKeyPem []byte, publicKeyPem []byte, err error) {
	var publicKeyWriter *bytes.Buffer = bytes.NewBufferString("")
	var privateKeyWriter *bytes.Buffer = bytes.NewBufferString("")
	defer func() {
		publicKeyWriter = nil
		privateKeyWriter = nil
	}()
	// 生成私钥文件
	privateKey, publicKey, err := CreateRsaKeys(keyLength)
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

func ParseX509CertificateFile(path string) (cert *x509.Certificate, err error) {
	certContent, err := os.ReadFile(path)
	if err != nil {
		err = fmt.Errorf("读取CA证书内容出错: %+v", err)
		return
	}

	return ParseX509CertificateContent(certContent)
}
func ParseX509CertificateContent(content []byte) (cert *x509.Certificate, err error) {
	blockCert, _ := pem.Decode(content)
	//fmt.Println("ParseCertificate b:", blockCert)

	cert, err = x509.ParseCertificate(blockCert.Bytes)
	//fmt.Println("ParseCertificate err:", err)
	if err != nil {
		err = fmt.Errorf("解析证书内容出错: %+v", err)
		return
	}

	return
}

func ParseX509PrivateKeyFile(path string) (priKey *rsa.PrivateKey, err error) {
	keyContent, err := os.ReadFile(path)
	if err != nil {
		err = fmt.Errorf("读取KEY内容出错: %+v", err)
		return
	}
	priKey, err = ParseX509PrivateKeyContent(keyContent)

	return
}
func ParseX509PrivateKeyContent(content []byte) (priKey *rsa.PrivateKey, err error) {
	blockKey, _ := pem.Decode(content)
	if nil == blockKey {
		err = fmt.Errorf("private key decode error")
		return
	}
	switch blockKey.Type {
	case "RSA PRIVATE KEY":
		priKey, err = x509.ParsePKCS1PrivateKey(blockKey.Bytes)
	// RFC5208 - https://tools.ietf.org/html/rfc5208
	case "PRIVATE KEY":
		tmp, err1 := x509.ParsePKCS8PrivateKey(blockKey.Bytes)
		if nil != tmp {
			if keyTmp, ok := tmp.(*rsa.PrivateKey); ok {
				priKey = keyTmp
			}
		}
		if nil != err1 {
			err = err1
		}

	default:
		return nil, fmt.Errorf("unsupported key type %q", blockKey.Type)
	}

	if err != nil {
		err = fmt.Errorf("解析KEY内容出错: %+v", err)
		return
	}

	return
}

func ParseX509PublicKeyFile(path string) (pubKey *rsa.PublicKey, err error) {
	keyContent, err := os.ReadFile(path)
	if err != nil {
		err = fmt.Errorf("读取KEY内容出错: %+v", err)
		return
	}
	pubKey, err = ParseX509PublicKeyContent(keyContent)

	return
}
func ParseX509PublicKeyContent(content []byte) (pubKey *rsa.PublicKey, err error) {
	block, _ := pem.Decode(content)
	if nil == block {
		err = fmt.Errorf("public key error")
		return
	}
	pkixPublicKey, err := x509.ParsePKIXPublicKey(block.Bytes)
	if nil != err {
		return
	}
	pubKey, ok := pkixPublicKey.(*rsa.PublicKey)
	if !ok {
		err = fmt.Errorf("public key type error")
		return
	}

	return
}
