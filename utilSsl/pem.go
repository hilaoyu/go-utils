package utilSsl

import (
	"encoding/pem"
	"os"
)

func ParsePemCertFile(path string) (pkey []byte, cert []byte) {
	b, _ := os.ReadFile(path)

	return ParsePemCertContent(b)
}

func ParsePemCertContent(content []byte) (pkey []byte, cert []byte) {
	var pemBlocks []*pem.Block
	var v *pem.Block
	pem.Decode(content)
	for {
		v, content = pem.Decode(content)
		if v == nil {
			break
		}
		if v.Type == "PRIVATE KEY" {
			pkey = pem.EncodeToMemory(v)
		} else {
			pemBlocks = append(pemBlocks, v)
		}
	}

	cert = pem.EncodeToMemory(pemBlocks[0])

	return
}
