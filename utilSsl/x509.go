package utilSsl

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"
)

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
	blockKey, _ := pem.Decode(keyContent)
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
func ParseX509PrivateKeyContent(content []byte) (priKey *rsa.PrivateKey, err error) {
	blockKey, _ := pem.Decode(content)
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
