package utils

import (
	"crypto/md5"
	"encoding/hex"
	"io"
)

//SHA1  加密
func Md5String(origin string) string {
	s := md5.New()
	io.WriteString(s, origin)
	return hex.EncodeToString(s.Sum(nil))
}

func Md5Byte(origin []byte) string {
	s := md5.New()
	s.Write(origin)
	return hex.EncodeToString(s.Sum(nil))
}

func Md516String(origin string) string {
	s := md5.New()
	io.WriteString(s, origin)
	b := s.Sum(nil)
	return hex.EncodeToString(b[4:12])
}
