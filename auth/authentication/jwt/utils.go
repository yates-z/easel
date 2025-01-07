package jwt

import (
	"encoding/base64"
	"strings"
)

// Base64URLEncode 执行 URL 安全的 Base64 编码
func Base64URLEncode(data []byte) string {
	return strings.TrimRight(base64.URLEncoding.EncodeToString(data), "=")
}

// Base64URLDecode 执行 URL 安全的 Base64 解码
func Base64URLDecode(encoded string) ([]byte, error) {
	if len(encoded)%4 != 0 {
		padding := 4 - len(encoded)%4
		encoded += strings.Repeat("=", padding)
	}
	return base64.URLEncoding.DecodeString(encoded)
}
