package utils

import (
	"crypto/hmac"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
)

// 计算hash1
func ComputeHash1(message string, secret string) string {
	h := hmac.New(sha1.New, []byte(secret))
	h.Write([]byte(message))
	// 转成十六进制
	return hex.EncodeToString(h.Sum(nil))
}

// 计算HmacSha256
func ComputeHmacSha256(message string, secret string) string {
	key := []byte(secret)
	h := hmac.New(sha256.New, key)
	h.Write([]byte(message))
	// 转成十六进制
	return hex.EncodeToString(h.Sum(nil))

}

// 编码Base64
func EncodeBase64(str string) string {
	return base64.StdEncoding.EncodeToString([]byte(str))
}
