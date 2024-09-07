package utils

import (
	"encoding/base64"
	"strings"
)

func Base64EncodeFormByte(bin []byte) []byte {
	e64 := base64.StdEncoding

	maxEncLen := e64.EncodedLen(len(bin))
	encBuf := make([]byte, maxEncLen)

	e64.Encode(encBuf, bin)
	return encBuf
}

func Base64EncodeUrlSafe(bin []byte) string {
	str := base64.StdEncoding.EncodeToString(bin)
	str = strings.ReplaceAll(str, "+", "-")
	str = strings.ReplaceAll(str, "/", "_")
	str = strings.ReplaceAll(str, "=", "")
	return str
}
