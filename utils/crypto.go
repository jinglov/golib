package utils

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"io"
)

//Md5  加密

func Md5(origin interface{}) (string, error) {
	s := md5.New()
	switch origin.(type) {
	case string:
		io.WriteString(s, origin.(string))
	case []byte:
		s.Write(origin.([]byte))
	default:
		return "", errors.New("not supper this type")
	}
	return hex.EncodeToString(s.Sum(nil)), nil
}
