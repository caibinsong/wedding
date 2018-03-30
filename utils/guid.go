package utils

import (
	"crypto/md5"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"io"
	"strings"
)

//生成32位md5字串
func GetMd5String(s string) string {
	h := md5.New()
	h.Write([]byte(s))
	return hex.EncodeToString(h.Sum(nil))
}

func GetGuid() string {
	b := make([]byte, 48)
	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		return ""
	}
	guid := strings.ToUpper(GetMd5String(base64.URLEncoding.EncodeToString(b)))
	return fmt.Sprintf("%s-%s-%s-%s-%s", guid[:8], guid[8:12], guid[12:16], guid[16:20], guid[20:32])
}

func GetRandStr() string {
	b := make([]byte, 48)
	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		return ""
	}
	return strings.ToLower(GetMd5String(base64.URLEncoding.EncodeToString(b)))
}
