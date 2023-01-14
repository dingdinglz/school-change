package tools

import (
	"crypto/md5"
	"encoding/hex"
)

// MD5 取字符串的md5
func MD5(src string) string {
	m := md5.New()
	m.Write([]byte(src))
	return hex.EncodeToString(m.Sum(nil))
}
